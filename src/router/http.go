package router

import (
	"context"
	"net/http"
	"services/initializer"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fzerorubigd/xhandler"
	"github.com/fzerorubigd/xmux"
)

type initRouter struct {
}

func (i initRouter) Initialize(ctx context.Context) {
	mux := xmux.New()
	mux.POST("/get", xhandler.HandlerFuncC(getAd))

	srv := &http.Server{Addr: listenAddress, Handler: xhandler.New(ctx, mux)}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Error(err)
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
	initializer.Register(&initRouter{})
}
