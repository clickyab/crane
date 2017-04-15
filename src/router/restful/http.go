package restful

import (
	"context"
	"net/http"
	"router/restful/internal/middlewares"
	"services/initializer"
	"time"

	"router/restful/internal/renderer"

	"net/url"

	"fmt"

	"services/assert"

	"github.com/Sirupsen/logrus"
	"github.com/fzerorubigd/xhandler"
	"github.com/fzerorubigd/xmux"
)

type ctxKey int

const rendererKey ctxKey = iota

type initRouter struct {
}

func (i initRouter) Initialize(ctx context.Context) {
	mux := xmux.New()
	mux.POST("/get/:key", xhandler.HandlerFuncC(middlewares.Recovery(middlewares.Logger(getAd))))

	pixel, err := url.Parse(fmt.Sprintf("http://%s/track", domain))
	assert.Nil(err)

	nCtx := context.WithValue(ctx, rendererKey, renderer.NewRestfulRenderer(pixel))
	srv := &http.Server{Addr: listenAddress, Handler: xhandler.New(nCtx, mux)}
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
	initializer.Register(&initRouter{})
}
