package metrics

import (
	"github.com/clickyab/services/config"
	"github.com/prometheus/client_golang/prometheus"
)

var port = config.RegisterString("crane.metric.port", "9700", "")

var (

	// Carrier of incoming impressions
	Carrier = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_request_carrier",
			Help: "Counter of request carrier",
		},
		[]string{"supplier", "carrier"},
	)

	// Size of incoming impressions
	Size = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_request_size",
			Help: "Histogram of request size",
		},
		[]string{"supplier", "size", "mode", "publisher", "type"},
	)

	// Filter reasons
	Filter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_filtered_reason",
			Help: "Counter of filter",
		},
		[]string{"supplier", "reason"},
	)

	Location = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_request_location",
			Help: "Counter for location",
		},
		[]string{"supplier", "latitude", "longitude", "country", "province", "isp", "hash"},
	)

	// Duration for getting response time
	Duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "demand_request_duration_millisecond",
			Help: "Histogram of request duration",
			Buckets: []float64{
				.0002,
				.0004,
				.0006,
				.0008,
				.001,
				.002,
				.004,
				.006,
				.008,
				.01,
				.02,
				.04,
				.06,
				.08,
				.1,
				.2,
				.4,
				.6,
				.8,
				1,
			},
		},
		[]string{"supplier", "route"},
	)

	// CounterRequest total request
	CounterRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_request_total",
			Help: "Total number of request",
		},
		[]string{"supplier", "route"},
	)
)
