package controller

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"services/assert"
	"services/config"
	"services/httplib/middlewares"
	"services/initializer"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/Sirupsen/logrus"
	echo "gopkg.in/labstack/echo.v3"
)

// Routes the base rote structure
type Routes interface {
	// Routes is for adding new controller
	Routes(r *echo.Echo, mountPoint string)
}

var (
	engine *echo.Echo
	all    []Routes

	cors = config.RegisterBoolean("services.httplib.controller.cors", true, "http controller cors")
	// this is development mode
	devel      = config.RegisterBoolean("core.devel_mode", true, "core developer mode")
	mountPoint = config.RegisterString("services.httplib.controller.mount_point", "/api", "http controller mount point")
	listen     *string
)

// Register a new controller class
func Register(c ...Routes) {
	all = append(all, c...)
}

// Why i end up with this fucking name?
type master struct {
}

// Initialize the controller
func (*master) Initialize(ctx context.Context) {
	engine = echo.New()
	mid := []echo.MiddlewareFunc{middlewares.Recovery, middlewares.Logger}
	if *cors {
		mid = append(mid, middlewares.CORS())
	}
	engine.Use(mid...)
	for i := range all {
		all[i].Routes(engine, *mountPoint)
	}

	//engine.SetLogLevel(log.DEBUG)
	if *devel {
		assetHandler := http.FileServer(rice.MustFindBox("../statics/swagger/").HTTPBox())
		engine.Any("/swagger/*", func(c echo.Context) error {
			http.StripPrefix("/swagger/", assetHandler).
				ServeHTTP(c.Response(), c.Request())
			return nil
		})
	}
	engine.Logger = NewLogger()

	go func() {
		if err := engine.Start(*listen); err != nil {
			logrus.Error(err)
		}
	}()

	done := ctx.Done()
	assert.NotNil(done, "[BUG] the done channel is nil")
	go func() {
		<-done
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		engine.Shutdown(ctx)
	}()
}

func init() {
	initializer.Register(&master{}, 100)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	listen = config.RegisterString("services.httplib.controller.listen", fmt.Sprintf(":%s", port), "http controller port")
}
