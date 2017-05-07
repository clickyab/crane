package broker

import (
	"services/assert"
	"services/config"
)

var (
	develMode = config.RegisterBoolean("core.devel_mode", false, "development mode")
	testMode  = config.RegisterBoolean("services.broker.test_mode", false, "test mode for development")
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

// Interface is the base broker interface in system
type Interface interface {
	// Publish is the async publisher for the broker
	Publish(Job)
}

var activeBroker Interface

// SetActiveBroker is a gateway to set active broker for this service
func SetActiveBroker(b Interface) {
	if *develMode && *testMode {
		return
	}
	assert.Nil(activeBroker, "[BUG] active broker is already set")
	activeBroker = b
}

// Publish try to Publish a job into system using the broker
func Publish(j Job) {
	assert.NotNil(activeBroker, "[BUG] active broker is not set")
	if *develMode && *testMode {
		return
	}
	activeBroker.Publish(j)
}
