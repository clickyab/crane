package metrics

import (
	"github.com/clickyab/services/config"
	"github.com/prometheus/client_golang/prometheus"
)

var port = config.RegisterString("crane.metric.port", "9700", "")

var (

	// Price of in and out
	Price = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_price",
			Help: "Counter of asset",
		},
		[]string{"price", "publisher", "io"},
	)

	// Asset of items
	Asset = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "supplier_asset",
			Help: "Counter of asset",
		},
		[]string{"list"},
	)

	// Impression of bid
	Impression = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_impression",
			Help: "Counter of impression",
		},
		[]string{"sup", "cid"},
	)

	// Loaded campaigns
	Loaded = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_loaded_campaign",
			Help: "loaded_ad",
		},
		[]string{"cid"},
	)

	// Click of incoming impressions
	Click = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_click",
			Help: "Counter of click",
		},
		[]string{"sup", "cid"},
	)

	// Campaigns of incoming impressions
	Campaigns = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_campaign",
			Help: "Counter of campaign",
		},
		[]string{"sup", "cid"},
	)

	// Carrier of incoming impressions
	Carrier = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_carrier",
			Help: "Counter of request carrier",
		},
		[]string{"sup", "carrier"},
	)

	// Size of incoming impressions
	Size = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_size",
			Help: "Histogram of request size",
		},
		[]string{"sup", "size", "io"},
	)

	// Filter reasons
	Filter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_filtered_reason",
			Help: "Counter of filter",
		},
		[]string{"sup", "reason"},
	)

	// Location of requests
	Location = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_location",
			Help: "Counter for location",
		},
		[]string{"sup", "hash"},
	)

	// Duration for getting response time
	Duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "demand_duration",
			Help: "Histogram of request duration",
			Buckets: []float64{
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
		[]string{"sup", "route"},
	)

	// CounterRequest total request
	CounterRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demand_total",
			Help: "Total number of request",
		},
		[]string{"sup", "route"},
	)
)
