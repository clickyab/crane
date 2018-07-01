package statics

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
	"github.com/rs/xhandler"
	"github.com/sirupsen/logrus"
)

type controller struct {
}

func (c controller) serveRobot(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	logrus.Warn("server robot route")
	logrus.Warn(os.Getenv("ROOT"))
	fullPath, _ := filepath.Abs(os.Getenv("ROOT") + "/statics/robots.txt")

	http.ServeFile(w, r, fullPath)
}

func (c controller) Routes(m framework.Mux) {
	m.RootMux().GET("/robots.txt", xhandler.HandlerFuncC(c.serveRobot))
}

func init() {
	router.Register(controller{})
}
