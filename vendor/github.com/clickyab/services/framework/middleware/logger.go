package middleware

import (
	"net/http"
	"time"

	"fmt"

	"context"

	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/xlog"
	"github.com/sirupsen/logrus"
)

type logger struct {
}

func (logger) Handler(next framework.Handler) framework.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

		// Start timer
		start := time.Now()
		// Process request
		wr := &wrapper{original: w}
		ray := <-random.ID
		ctx = xlog.SetFields(
			ctx,
			logrus.Fields{
				"cy-ray": ray,
			},
		)

		defer func() {
			latency := time.Since(start)
			xlog.GetWithFields(
				ctx,
				logrus.Fields{
					"domain":    r.Host,
					"method":    r.Method,
					"path":      r.URL.Path,
					"latency":   fmt.Sprint(latency),
					"status":    wr.status,
					"len":       wr.total,
					"client_ip": framework.RealIP(r),
				},
			).Debug(http.StatusText(wr.status))
		}()

		wr.Header().Add("cy-ray", ray)
		next(ctx, wr, r)
	}
}

func (logger) PreRoute() bool {
	return true
}

type wrapper struct {
	original http.ResponseWriter
	total    int
	status   int
}

func (w *wrapper) Header() http.Header {
	return w.original.Header()
}

func (w *wrapper) Write(b []byte) (int, error) {
	x, err := w.original.Write(b)
	w.total += x
	return x, err
}

func (w *wrapper) WriteHeader(c int) {
	w.status = c
	w.original.WriteHeader(c)
}

// Logger is the middleware for log system
func Logger() framework.GlobalMiddleware {
	return &logger{}
}
