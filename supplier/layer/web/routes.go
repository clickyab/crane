package web

import (
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

type controller struct {
}

func (controller) Routes(r framework.Mux) {
	r.GET("multi-js", "/multi.js", js)
	// fix the following line after writing the actual route
	r.GET("multi-ad", "/todo/fix/me", js)
}

func init() {
	router.Register(&controller{})
}
