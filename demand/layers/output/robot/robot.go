package robot

import (
	"context"
	"net/http"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
	"github.com/rs/xhandler"
)

var robotTxt = `User-agent: *
Disallow: /api/
Disallow: /api/*`

type controller struct {
}

func (c controller) serveRobot(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	assert.Nil(w.Write([]byte(robotTxt)))
}

func (c controller) Routes(m framework.Mux) {
	m.RootMux().GET("/robots.txt", xhandler.HandlerFuncC(c.serveRobot))
}

func init() {
	router.Register(controller{})
}
