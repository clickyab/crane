package middlewares

import (
	"net/http"
	"time"

	"context"

	"clickyab.com/exchange/services/httplib"

	"github.com/Sirupsen/logrus"
	"github.com/fzerorubigd/xhandler"
)

type dummyWriter struct {
	w      http.ResponseWriter
	status int
}

func (dw *dummyWriter) Header() http.Header {
	return dw.w.Header()
}

func (dw *dummyWriter) Write(b []byte) (int, error) {
	return dw.w.Write(b)
}

func (dw *dummyWriter) WriteHeader(s int) {
	dw.status = s
	dw.w.WriteHeader(s)
}

// Logger is the middleware for log system
func Logger(next xhandler.HandlerFuncC) xhandler.HandlerFuncC {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		// Start timer
		start := time.Now()
		ip := httplib.RealIP(r)
		// Process request
		logrus.WithFields(
			logrus.Fields{
				"Method":   r.Method,
				"Path":     r.URL.Path,
				"ClientIP": ip,
			},
		).Debug("STARTED")
		dummy := &dummyWriter{w: w, status: http.StatusOK}
		next(c, dummy, r)
		latency := time.Since(start)
		logrus.WithFields(
			logrus.Fields{
				"Method":   r.Method,
				"Path":     r.URL.Path,
				"Latency":  latency,
				"ClientIP": ip,
				"Status":   dummy.status,
			},
		).Debug("DONE WITH: " + http.StatusText(dummy.status))
	}
}
