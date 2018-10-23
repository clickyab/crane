package metrics

import (
	"context"
	"net/http"

	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/safe"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type initer struct {
}

func (initer) Initialize(ctx context.Context) {
	prometheus.MustRegister(Duration)
	prometheus.MustRegister(CounterRequest)
	prometheus.MustRegister(Size)
	prometheus.MustRegister(Filter)

	http.Handle("/metrics", promhttp.Handler())
	safe.GoRoutine(ctx, func() {
		http.ListenAndServe(":"+port.String(), nil)

	})
}

func init() {
	initializer.Register(&initer{}, 100)
}
