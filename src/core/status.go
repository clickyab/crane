package core

import (
	"context"
	"net/http"

	"github.com/fzerorubigd/xmux"
)

func (p providerData) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	p.provider.Status(ctx, w, r)
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
