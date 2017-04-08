package core

import (
	"assert"
	"context"
	"entity"
	"net/http"
	"sync"
	"time"

	"github.com/rs/xmux"
	oldctx "golang.org/x/net/context"
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

func (p providerData) ServeHTTPC(ctx oldctx.Context, w http.ResponseWriter, r *http.Request) {
	p.provider.Status(ctx, w, r)
}

func (p *providerData) watch(ctx context.Context, imp entity.Impression) map[string]entity.Advertise {
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

	allProviders[provider.Name()] = providerData{
		name:     provider.Name(),
		provider: provider,
		timeout:  timeout,
	}
}

// Reset remove all providers
func Reset() {
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
		if limit == 0 {
			limit = imp.Source().Supplier().FloorCPM()
		}
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
