package show

import (
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

type controller struct {
}

func (controller) Routes(m framework.Mux) {
	m.GET("banner", showPath, showBanner)
}

func init() {
	router.Register(controller{})
}
