package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/middleware"
	"github.com/clickyab/services/initializer"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	"github.com/sirupsen/logrus"
	onion "gopkg.in/fzerorubigd/onion.v3"
)

var (
	engine *xmux.Mux
	all    = []framework.Routes{}
	mid    = []framework.GlobalMiddleware{middleware.Logger()}

	// this is development mode
	mountPoint = config.RegisterString("services.framework.controller.mount_point", "/api", "http controller mount point")
	hammerTime = config.RegisterDuration("services.framework.controller.graceful_wait", 100*time.Millisecond, "the time for framework to stop for graceful shutdown")
	listen     onion.String
)

type fake struct {
	base framework.Handler
}

func (f fake) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	f.base(ctx, w, r)
}

type initer struct {
}

func (i *initer) Initialize(ctx context.Context) {
	engine = xmux.New()

	var (
		pre  []framework.GlobalMiddleware
		post []framework.GlobalMiddleware
	)
	for i := range mid {
		if mid[i].PreRoute() {
			pre = append(pre, mid[i])
		} else {
			post = append(post, mid[i])
		}
	}

	fPre := func(next framework.Handler) framework.Handler {
		for i := range pre {
			next = pre[i].Handler(next)
		}

		return next
	}

	fPost := func(next framework.Handler) framework.Handler {
		for i := range post {
			next = post[i].Handler(next)
		}

		return next
	}

	xm := &xmuxer{
		root:       engine,
		middleware: fPost,
	}
	mp := mountPoint.String()
	if mp != "" {
		xm.group = engine.NewGroup(mp)
	} else {
		xm.engine = engine
	}

	for i := range all {
		all[i].Routes(xm)
	}
	// Append some generic middleware, to handle recovery and log
	handler := middleware.Recovery(
		xhandler.New(context.Background(), fake{base: fPre(engine.ServeHTTPC)}).ServeHTTP,
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
		ctx, _ := context.WithTimeout(context.Background(), hammerTime.Duration())
		server.Shutdown(ctx)
	}()
}

// Register a new controller class
func Register(c ...framework.Routes) {
	all = append(all, c...)
}

// RegisterGlobalMiddleware is a function to register a global middleware
func RegisterGlobalMiddleware(g framework.GlobalMiddleware) {
	mid = append(mid, g)
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
