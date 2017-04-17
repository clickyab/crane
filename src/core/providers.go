package core

import (
	"context"
	"entity"
	"net/http"
	"services/assert"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fzerorubigd/xmux"
)

var (
	allProviders = make(map[string]providerData)
	lock         = &sync.RWMutex{}
)

type providerData struct {
	name     string
	provider entity.Demand
	timeout  time.Duration
}

func (p providerData) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	p.provider.Status(ctx, w, r)
}

func (p *providerData) watch(ctx context.Context, imp entity.Impression) map[string]entity.Advertise {
	logrus.Debugf("Watch in for %s", p.provider.Name())
	defer logrus.Debugf("Watch out for %s", p.provider.Name())
	done := ctx.Done()
	assert.NotNil(done)

	res := make(map[string]entity.Advertise)
	// the cancel is not required here. the parent is the hammer :)
	rCtx, _ := context.WithTimeout(ctx, p.timeout)
	chn := make(chan map[string]entity.Advertise, len(imp.Slots()))
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
func Register(provider entity.Demand, timeout time.Duration) {
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

// Mount is called to create status page for all providers
func Mount(m *xmux.Mux) {
	lock.RLock()
	defer lock.RUnlock()

	for i := range allProviders {
		m.GET("/"+allProviders[i].name, allProviders[i])
		m.POST("/"+allProviders[i].name, allProviders[i])
		m.OPTIONS("/"+allProviders[i].name, allProviders[i])
		m.DELETE("/"+allProviders[i].name, allProviders[i])
		m.HEAD("/"+allProviders[i].name, allProviders[i])
	}
}

// Call is for getting the current ads for this imp
func Call(ctx context.Context, imp entity.Impression) map[string][]entity.Advertise {
	rCtx, cnl := context.WithTimeout(ctx, maximumTimeout)
	defer cnl()

	wg := sync.WaitGroup{}
	l := len(allProviders)
	wg.Add(l)
	allRes := make(chan map[string]entity.Advertise, l)
	lock.RLock()
	for i := range allProviders {
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

	res := make(map[string][]entity.Advertise)
	for provided := range allRes {
		for j := range provided {
			if provided[j].MaxCPM() > limit {
				res[j] = append(res[j], provided[j])
			}
		}
	}

	return res
}
