package router

import (
	"context"
	"net/http"
	"time"

	"os"

	"fmt"

	rice "github.com/GeertJohan/go.rice"
	"github.com/Sirupsen/logrus"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/middleware"
	"github.com/clickyab/services/initializer"
	"github.com/rs/cors"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	onion "gopkg.in/fzerorubigd/onion.v3"
)

var (
	engine *xmux.Mux
	all    []Routes

	// this is development mode
	mountPoint = config.RegisterString("services.framework.controller.mount_point", "/api", "http controller mount point")
	swagger    = config.RegisterBoolean("services.framework.swagger", false, "is any swagger code available?")
	listen     onion.String
)

// Routes the base rote structure
type Routes interface {
	// Routes is for adding new controller
	Routes(r *xmux.Mux, mountPoint string)
}

type initer struct {
}

func (i *initer) Initialize(ctx context.Context) {
	engine = xmux.New()
	c := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true // TODO : write the real code here
		},
	})
	//engine.SetLogLevel(log.DEBUG)
	if swagger.Bool() {
		assetHandler := http.FileServer(rice.MustFindBox("../statics/swagger/").HTTPBox())
		framework.Any(engine, "/swagger/*", func(_ context.Context, w http.ResponseWriter, r *http.Request) {
			http.StripPrefix("/swagger/", assetHandler).ServeHTTP(w, r)
		})
	}

	for i := range all {
		all[i].Routes(engine, mountPoint.String())
	}

	// Append some generic middleware, to handle recovery, log and CORS
	handler := middleware.Recovery(
		middleware.Logger(
			c.Handler(xhandler.New(context.Background(), engine)).ServeHTTP,
		),
	)
	server := &http.Server{Addr: listen.String(), Handler: handler}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.Error(err)
		}
	}()

	done := ctx.Done()
	assert.NotNil(done, "[BUG] the done channel is nil")
	go func() {
		<-done
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		server.Shutdown(ctx)
	}()
}

// Register a new controller class
func Register(c ...Routes) {
	all = append(all, c...)
}

func init() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	listen = config.RegisterString(
		"services.framework.listen",
		fmt.Sprintf(":%s", port),
		"address to listen for framework",
	)

	initializer.Register(&initer{}, 100)
}
