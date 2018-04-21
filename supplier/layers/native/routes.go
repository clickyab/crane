package native

import (
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

type controller struct {
}

// Routes is for registering routes
func (controller) Routes(r framework.Mux) {
	r.GET("nativejs", "/api/native.js", getNativeJS)
	r.GET("native", "/api/get/native", getNative)
}

func init() {
	router.Register(&controller{})
}
