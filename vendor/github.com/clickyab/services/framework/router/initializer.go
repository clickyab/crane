package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
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
