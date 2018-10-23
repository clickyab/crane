package metrics

import (
	"github.com/clickyab/services/config"
	"github.com/prometheus/client_golang/prometheus"
)

var port = config.RegisterString("crane.metric.port", "9700", "")

var (

	// Publisher counter
	Publisher = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_publisher",
			Help: "Counter for publishers",
		},
		[]string{"supplier", "publisher", "type"},
	)

	// Size of incoming impressions
	Size = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_request_size",
			Help: "Histogram of request size",
		},
		[]string{"supplier", "size", "mode"},
	)

	// Filter reasons
	Filter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_filtered_reason",
			Help: "Counter of filter",
		},
		[]string{"supplier", "reason"},
	)

	// Duration for getting response time
	Duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "demand_request_duration_millisecond",
			Help: "Histogram of request duration",
			Buckets: []float64{
				.0000,
				.0001,
				.0002,
				.0003,
				.0004,
				.0005,
				.0006,
				.0007,
				.0008,
				.0009,
				.001,
				.002,
				.003,
				.004,
				.005,
				.006,
				.007,
				.008,
				.009,
				.01,
				.02,
				.03,
				.04,
				.05,
				.06,
				.07,
				.08,
				.09,
				.1,
				.2,
				.3,
				.4,
				.5,
				.6,
				.7,
				.8,
				.9,
				1,
				2,
				3,
				4,
				5,
				6,
				7,
				8,
				9,
			},
		},
		[]string{"status", "supplier", "route"},
	)

	// CounterRequest total request
	CounterRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_request_total",
			Help: "Total number of request",
		},
		[]string{"status", "supplier", "route"},
	)
)
