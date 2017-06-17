package broker

import (
	"sync"

	"github.com/clickyab/services/assert"
)

// Job is a normal job
type Job interface {
	// Encode this job into string
	Encode() ([]byte, error)
	// The length of encoded data
	Length() int
	// Topic return the topic that this job is going to Publish into
	Topic() string
	// Key is partitioning key if this is possible for the broker
	Key() string
	// Report is called on every single message on error or success. if the error is nil, then the
	// broker handle it successfully.
	Report() func(error)
}

// Delivery is the job to consumer
type Delivery interface {
	Decode(v interface{}) error
	// Ack delegates an acknowledgement through the Acknowledger interface that the client or server has finished work on a delivery.
	Ack(multiple bool) error
	// Nack negatively acknowledge the delivery of message(s) identified by the delivery tag from either the client or server.
	Nack(multiple, requeue bool) error
	// Reject delegates a negatively acknowledgement through the Acknowledger interface.
	Reject(requeue bool) error
}

// Consumer is the side the workers on it
type Consumer interface {
	// Topic return the topic that this worker want to listen to it
	Topic() string
	// Queue is the queue that this want to listen to
	Queue() string
	// Consume return a channel to put jobs into
	Consume() chan<- Delivery
}

// Publisher is the base broker interface in system
type Publisher interface {
	// Publish is the async publisher for the broker
	Publish(Job)
}

// Interface is the full broker interface
type Interface interface {
	Publisher

	// RegisterConsumer try to register a consumer in system
	RegisterConsumer(Consumer) error
}

var (
	activeBroker Publisher
	lock         = sync.RWMutex{}
)

// SetActiveBroker is a gateway to set active broker for this service
func SetActiveBroker(b Publisher) {
	lock.Lock()
	defer lock.Unlock()
	assert.Nil(activeBroker, "[BUG] active broker is already set")
	activeBroker = b
}

// Publish try to Publish a job into system using the broker
func Publish(j Job) {
	lock.RLock()
	defer lock.RUnlock()

	assert.NotNil(activeBroker, "[BUG] active broker is not set")
	activeBroker.Publish(j)
}

// RegisterConsumer is the endpoint to register a consumer in active broker
func RegisterConsumer(consumer Consumer) error {
	lock.RLock()
	defer lock.RUnlock()

	assert.NotNil(activeBroker, "[BUG] active broker is not set")
	// there is no need to check if the broker support for consumer.
	// go do a panic here if its not and its ok
	return activeBroker.(Interface).RegisterConsumer(consumer)
}
