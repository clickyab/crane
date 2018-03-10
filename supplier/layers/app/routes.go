package app

import (
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

type controller struct {
}

// Routes is for registering routes
func (controller) Routes(r framework.Mux) {
	r.GET("app-single-ad", "/api/getapp", getApp)
	r.GET("app-single-ad-bc", "/ads/inapp.php", getApp)
	r.GET("app-json-app", "/ads/json-inapp.php", getInappJSON)
}

func init() {
	router.Register(&controller{})
}
