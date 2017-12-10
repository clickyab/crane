package click

import (
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

type controller struct {
}

func (controller) Routes(m framework.Mux) {
	m.GET("banner", clickPath, clickBanner)
}

func init() {
	router.Register(controller{})
}
