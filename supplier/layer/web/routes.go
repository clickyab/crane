package web

import (
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

type controller struct {
}

func (controller) Routes(r framework.Mux) {
	r.GET("multi-js", "/api/multi.js", jsV2)
	r.GET("show-js", "/show.js", showV1.ServeHTTPC)
	// fix the following line after writing the actual route
	r.GET("multi-ad", "/api/getad", getAd)

}

func init() {
	router.Register(&controller{})
}
