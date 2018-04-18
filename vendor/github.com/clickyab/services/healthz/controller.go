package healthz

import (
	"context"
	"net/http"

	"fmt"

	"time"

	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
	"github.com/rs/xhandler"
	"github.com/sirupsen/logrus"
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

	w.Header().Set("time", time.Now().String())
	if len(errs) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		for i := range errs {
			fmt.Fprint(w, errs[i].Error())
		}
		return
	}
	framework.JSON(w, http.StatusOK, struct {
		Time string `json:"time"`
	}{
		Time: time.Now().String(),
	})
}

func (r route) Routes(mux framework.Mux) {
	mux.RootMux().GET("/healthz", xhandler.HandlerFuncC(r.check))
}

func init() {
	router.Register(&route{})
}
