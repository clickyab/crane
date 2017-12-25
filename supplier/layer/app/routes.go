package app

import (
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

type controller struct {
}

// Routes is for registering routes
func (controller) Routes(r framework.Mux) {
	r.GET("app-single-ad", "/getapp", getApp)
}

func init() {
	router.Register(&controller{})
}
