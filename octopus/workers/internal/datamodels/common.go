package datamodels

import "sync"

var (
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

	RequestInCount     int64 //demand
	RequestOutCount    int64 //demand
	ImpressionInCount  int64 //imp,demand
	ImpressionOutCount int64 //demand,win
	WinCount           int64 //win
	WinBid             int64 //win
	DeliverCount       int64 //show
	DeliverBid         int64 //show
	Profit             int64 //show

	Acknowledger Acknowledger
	WorkerID     string
}

// Aggregator is a helper type to handle the entire process, and hey, its mock-able!
type Aggregator interface {
	Channel() chan<- TableModel
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
