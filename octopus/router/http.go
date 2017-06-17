// Package router is a glu package to mix all parts together
package router

import (
	"context"
	"net/http"
	"time"

	"clickyab.com/exchange/octopus/router/internal/demands"
	"clickyab.com/exchange/octopus/router/internal/middlewares"
	"clickyab.com/exchange/octopus/router/internal/restful"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"

	"github.com/Sirupsen/logrus"
	"github.com/fzerorubigd/xhandler"
	"github.com/fzerorubigd/xmux"
)

var (
	listenAddress = config.RegisterString("exchange.router.listen", ":80", "exchange router listen address")
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
	mux.GET("/demands/status/:name", wrap(demands.Status))
	mux.POST("/demands/status/:name", wrap(demands.Status))
	mux.DELETE("/demands/status/:name", wrap(demands.Status))
	mux.HEAD("/demands/status/:name", wrap(demands.Status))
	mux.PUT("/demands/status/:name", wrap(demands.Status))
	mux.OPTIONS("/demands/status/:name", wrap(demands.Status))

	srv := &http.Server{Addr: *listenAddress, Handler: xhandler.New(ctx, mux)}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Debug(err)
		}
	}()
	logrus.Debugf("Server started on %s", *listenAddress)
	go func() {
		done := ctx.Done()
		if done != nil {
			<-done
			s, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
			srv.Shutdown(s)
			logrus.Debug("Server stopped")
		}
	}()

}

func init() {
	initializer.Register(&initRouter{}, 100)
}
