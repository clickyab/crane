package restful

import (
	"context"
	"net/http"
	"octopus/router/restful/internal/middlewares"
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

func (i initRouter) Initialize(ctx context.Context) {
	mux := xmux.New()
	mux.POST("/get/:key", xhandler.HandlerFuncC(middlewares.Recovery(middlewares.Logger(getAd))))
	mux.GET("/pixel/:demand/:trackID", xhandler.HandlerFuncC(middlewares.Recovery(middlewares.Logger(trackPixel))))

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
