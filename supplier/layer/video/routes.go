package video

import (
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

type controller struct {
}

func (controller) Routes(r framework.Mux) {
	r.GET("vast", "/api/vast", vast)

}

func init() {
	router.Register(&controller{})
}
