package statics

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
	"github.com/rs/xhandler"
)

type controller struct {
}

func (c controller) serveRobot(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	fullPath, _ := filepath.Abs(os.Getenv("ROOT") + "/statics/robots.txt")

	http.ServeFile(w, r, fullPath)
}

func (c controller) Routes(m framework.Mux) {
	m.RootMux().GET("/robots.txt", xhandler.HandlerFuncC(c.serveRobot))
}

func init() {
	router.Register(controller{})
}
