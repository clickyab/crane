// Package router is a glu package to mix all parts together
package router

import (
	"context"
	"net/http"
	"octopus/router/internal/demands"
	"octopus/router/internal/middlewares"
	"octopus/router/internal/restful"
	"services/config"
	"services/initializer"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fzerorubigd/xhandler"
	"github.com/fzerorubigd/xmux"
)

var (
	listenAddress = config.RegisterString("exchange.router.listen", ":80", "exchnage router listen address")
)

type initRouter struct {
}

func wrap(in xhandler.HandlerFuncC) xhandler.HandlerFuncC {
	return xhandler.HandlerFuncC(middlewares.Recovery(middlewares.Logger(in)))
}

func (i initRouter) Initialize(ctx context.Context) {
	mux := xmux.New()
	// Restful route
	mux.POST("/rest/get/:key", wrap(restful.GetAd))
	mux.GET("/pixel/:demand/:trackID", wrap(restful.TrackPixel))
	// The demand status routes
	mux.GET("/demands/status/:name", wrap(demands.DemandStatus))
	mux.POST("/demands/status/:name", wrap(demands.DemandStatus))
	mux.DELETE("/demands/status/:name", wrap(demands.DemandStatus))
	mux.HEAD("/demands/status/:name", wrap(demands.DemandStatus))
	mux.PUT("/demands/status/:name", wrap(demands.DemandStatus))
	mux.OPTIONS("/demands/status/:name", wrap(demands.DemandStatus))

	srv := &http.Server{Addr: *listenAddress, Handler: xhandler.New(ctx, mux)}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Debug(err)
		}
	}()

	go func() {
		done := ctx.Done()
		if done != nil {
			<-done
			s, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
			srv.Shutdown(s)
		}
	}()

}

func init() {
	initializer.Register(&initRouter{}, 100)
}
