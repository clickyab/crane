package middleware

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/clickyab/services/framework"
)

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
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Start timer
		start := time.Now()
		// Process request
		wr := &wrapper{original: w}
		next(wr, r)

		latency := time.Since(start)
		logrus.WithFields(
			logrus.Fields{
				"Method":   r.Method,
				"Path":     r.URL.Path,
				"Latency":  latency,
				"ClientIP": framework.RealIP(r),
				"Status":   wr.status,
				"Len":      wr.total,
			},
		).Debug(http.StatusText(wr.status))
	}
}
