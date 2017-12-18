package web

import (
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

type controller struct {
}

func (controller) Routes(framework.Mux) {
	panic("implement me")
}

func init() {
	router.Register(&controller{})
}
