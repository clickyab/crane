package core

import (
	"context"
	"crane/entity"
	"net/http"
	"sync"
	"time"

	"services/assert"

	"github.com/fzerorubigd/xmux"
)

var allProviders []providerData

type providerData struct {
	name     string
	provider entity.Demand
	timeout  time.Duration
}

func (p providerData) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
func Register(name string, provider entity.Demand, timeout time.Duration) {
	for i := range allProviders {
		assert.True(allProviders[i].name != name, "[BUG] same name registered twice")
	}

	allProviders = append(
		allProviders,
		providerData{
			name:     name,
			provider: provider,
			timeout:  timeout,
		},
	)
}

// Mount is called to create status page for all providers
func Mount(m *xmux.Mux) {
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
	for i := range allProviders {
		go func(inner int) {
			defer wg.Done()
			res := allProviders[inner].watch(rCtx, imp)
			if res != nil {
				allRes <- res
			}
		}(i)
	}

	wg.Wait()
	// The close is essential here.
	close(allRes)
	res := make(map[string][]entity.Advertise)
	for provided := range allRes {
		for j := range provided {
			res[j] = append(res[j], provided[j])
		}
	}

	return res
}
