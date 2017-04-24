package core

import (
	"context"
	"octopus/exchange"
	"services/assert"
	"sync"
	"sync/atomic"
	"time"

	"errors"

	"github.com/Sirupsen/logrus"
)

var (
	allProviders = make(map[string]providerData)
	lock         = &sync.RWMutex{}
)

type providerData struct {
	name            string
	provider        exchange.Demand
	timeout         time.Duration
	callRateTracker int64
}

func (p *providerData) Skip() bool {
	x := atomic.AddInt64(&p.callRateTracker, 1)
	if x%int64(100/p.provider.CallRate()) != 0 {
		return false
	}
	return true
}

func (p *providerData) watch(ctx context.Context, imp exchange.Impression) map[string]exchange.Advertise {
	logrus.Debugf("Watch in for %s", p.provider.Name())
	defer logrus.Debugf("Watch out for %s", p.provider.Name())
	done := ctx.Done()
	assert.NotNil(done)

	res := make(map[string]exchange.Advertise)
	// the cancel is not required here. the parent is the hammer :)
	rCtx, _ := context.WithTimeout(ctx, p.timeout)
	chn := make(chan map[string]exchange.Advertise, len(imp.Slots()))
	go p.provider.Provide(rCtx, imp, chn)
	for {
		select {
		case <-done:
			// request is canceled
			return res
		case data, open := <-chn:
			for i := range data {
				res[i] = data[i]
			}
			if !open {
				return res
			}
		}
	}
}

// Register is used to handle new layer in system
func Register(provider exchange.Demand, timeout time.Duration) {
	lock.Lock()
	defer lock.Unlock()
	name := provider.Name()

	_, ok := allProviders[name]
	assert.False(ok, "[BUG] provider is already registered")

	allProviders[name] = providerData{
		name:     name,
		provider: provider,
		timeout:  timeout,
	}
}

// ResetProviders remove all providers
func ResetProviders() {
	lock.Lock()
	defer lock.Unlock()

	allProviders = make(map[string]providerData)
}

// Call is for getting the current ads for this imp
func Call(ctx context.Context, imp exchange.Impression) map[string][]exchange.Advertise {
	rCtx, cnl := context.WithTimeout(ctx, maximumTimeout)
	defer cnl()

	wg := sync.WaitGroup{}
	l := len(allProviders)
	wg.Add(l)
	allRes := make(chan map[string]exchange.Advertise, l)
	lock.RLock()
	for i := range allProviders {
		if !demandIsAllowed(imp, allProviders[i]) {
			continue
		}
		go func(inner string) {
			defer wg.Done()
			p := allProviders[inner]
			res := p.watch(rCtx, imp)
			if res != nil {
				allRes <- res
			}
		}(i)
	}
	lock.RUnlock()

	wg.Wait()
	// The close is essential here.
	close(allRes)
	var limit int64
	if !imp.UnderFloor() {
		limit = imp.Source().FloorCPM()
	}

	res := make(map[string][]exchange.Advertise)
	for provided := range allRes {
		for j := range provided {
			if provided[j].MaxCPM() > limit {
				res[j] = append(res[j], provided[j])
			}
		}
	}

	return res
}

func demandIsAllowed(impression exchange.Impression, data providerData) bool {
	if impression.Source().Name() == data.name {
		return false
	}
	if contains(data.provider.ExcludedSuppliers(), impression.Source().Name()) {
		return false
	}
	//if contains(data.provider.WhiteListCountries(),impression.Location().Country().Name) {
	//	return false
	//}
	//if contains(impression.Source().Supplier().ExcludedDemands(), data.name) {
	//	return false
	//}
	return true
}

func contains(s []string, t string) bool {
	if len(s) == 0 {
		return false
	}
	for _, a := range s {
		if a == t {
			return false
		}
	}
	return true
}

// GetDemand return demand by its name
func GetDemand(name string) (exchange.Demand, error) {
	lock.RLock()
	defer lock.RUnlock()
	val, found := allProviders[name]
	if !found {
		return nil, errors.New("demand not found")
	}
	return val.provider, nil
}
