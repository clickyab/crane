package ortb

import (
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

type controller struct {
}

func (controller) Routes(m framework.Mux) {
	m.POST("ortb", demandPath, openRTBInput)
}

func init() {
	router.Register(controller{})
}
