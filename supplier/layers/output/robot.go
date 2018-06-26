package output

import (
	"context"
	"net/http"

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
	_, err := w.Write([]byte(robotTxt))
	if err != nil {
		w.Write([]byte("error occurred"))
	}
}

func (c controller) Routes(m framework.Mux) {
	m.RootMux().GET("/robots.txt", xhandler.HandlerFuncC(c.serveRobot))
}

func init() {
	router.Register(controller{})
}
