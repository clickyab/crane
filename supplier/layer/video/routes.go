package video

import (
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

type controller struct {
}

func (controller) Routes(r framework.Mux) {
	r.GET("vast", "/ads/vast", vast)
	r.GET("jwplayer", "/api/vast.js", getJwplayer)
	r.GET("videojs", "/api/videojs.js", getVideojs)
}

func init() {
	router.Register(&controller{})
}
