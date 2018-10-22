package metrics

import (
	"github.com/clickyab/services/config"
	"github.com/prometheus/client_golang/prometheus"
)

var port = config.RegisterString("crane.metric.port", "9700", "")

var (
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
	// CounterImpression for total impression ( this should be more than total request )
	CounterImpression = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_request_impression",
			Help: "Total number of impressions",
		},
		[]string{"status", "supplier", "route"},
	)

	// CounterSeat is for our response ( must be less then impression )
	CounterSeat = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_request_seat",
			Help: "Total number of seats",
		},
		[]string{"status", "supplier", "route"},
	)
)
