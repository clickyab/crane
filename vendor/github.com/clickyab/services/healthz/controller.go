package healthz

import (
	"context"
	"net/http"

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/clickyab/services/framework/router"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
)

type route struct {
}

func (r route) check(ctx context.Context, w http.ResponseWriter, rq *http.Request) {
	lock.RLock()
	defer lock.RUnlock()

	var (
		errs []error
	)

	for i := range all {
		if err := all[i].Healthy(ctx); err != nil {
			logrus.Error(err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		for i := range errs {
			fmt.Fprint(w, errs[i].Error())
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (r route) Routes(mux *xmux.Mux, mountPoint string) {
	mux.GET("/healthz", xhandler.HandlerFuncC(r.check))
}

func init() {
	router.Register(&route{})
}
