package datamodels

import (
	"sync"
	"time"

	"clickyab.com/exchange/services/assert"
)

var (
	epoch     time.Time
	singleton Aggregator
	lock      sync.RWMutex
)

// Acknowledger is the job to consumer, good parts only
type Acknowledger interface {
	// Ack delegates an acknowledgement through the Acknowledger interface that the client or server has finished work on a delivery.
	Ack(multiple bool) error
	// Nack negatively acknowledge the delivery of message(s) identified by the delivery tag from either the client or server.
	Nack(multiple, requeue bool) error
	// Reject delegates a negatively acknowledgement through the Acknowledger interface.
	Reject(requeue bool) error
}

// TableModel is the model for counting data and aggregate them into on query
type TableModel struct {
	Supplier string // All
	Source   string // All
	Demand   string // All
	Time     int64  // All

	ImpressionRequest int64 // Impression
	ImpressionSlots   int64 // Impression
	DemandRequest     int64 // Demand
	DemandAds         int64 // Demand
	WinCount          int64 // Winner
	ShowCount         int64 // Winner
	ImpressionBid     int64 // Demand
	WinBid            int64 // Winner
	ShowBid           int64 // Show

	Acknowledger Acknowledger
}

// Aggregator is a helper type to handle the entire process, and hey, its mock-able!
type Aggregator interface {
	Channel() chan<- TableModel
}

// FactTableID is a helper function to get the fact table id from time
func FactTableID(tm time.Time) int64 {
	return int64(tm.Sub(epoch).Hours()) + 1
}

// RegisterAggregator to register an aggregator
func RegisterAggregator(a Aggregator) {
	lock.Lock()
	defer lock.Unlock()

	singleton = a
}

// ActiveAggregator return the current aggregator
func ActiveAggregator() Aggregator {
	lock.RLock()
	defer lock.RUnlock()

	return singleton
}

func init() {
	layout := "2006-01-02T15:04:05.000Z"
	str := "2017-03-21T00:00:00.000Z"
	var err error
	epoch, err = time.Parse(layout, str)
	assert.Nil(err)
}
