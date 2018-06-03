package statics

import (
	"context"
	"net/http"
	"path/filepath"

	"clickyab.com/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
	"github.com/rs/xhandler"
)

var (
	staticsPath = config.RegisterString("crane.statics.path", "/app/statics", "determine statics full path")
)

type controller struct {
}

// router handler to serve static files
func (c controller) serveStatic(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	fullPath, _ := filepath.Abs(staticsPath.String() + "/statics/" + r.URL.Path[len("/static/"):])

	http.ServeFile(w, r, fullPath)
}

func (c controller) serveRobot(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	fullPath, _ := filepath.Abs(staticsPath.String() + "/statics/robots.txt")

	http.ServeFile(w, r, fullPath)
}

func (c controller) Routes(m framework.Mux) {
	m.RootMux().GET("/robots.txt", xhandler.HandlerFuncC(c.serveRobot))
	m.RootMux().GET("/static/*path", xhandler.HandlerFuncC(c.serveStatic))
}

func init() {
	router.Register(controller{})
}
