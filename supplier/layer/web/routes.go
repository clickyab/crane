package web

import (
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

type controller struct {
}

func (controller) Routes(r framework.Mux) {
	r.GET("multi-js", "/multi.js", jsV2)
	r.RootMux().GET("/show.js", showV1)
	// fix the following line after writing the actual route
	r.GET("multi-ad", "/getad", getAd)

}

func init() {
	router.Register(&controller{})
}
