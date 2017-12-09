package kv

import (
	"time"

	"github.com/clickyab/services/assert"
)

// OneTimeSet is a simple way to set a one time ting, also any new call will
// not effect the value
type OneTimeSet interface {
	// Key return the current key
	Key() string
	// Set try to set a value, if its set for the first time,
	// then return its value
	// if not, return the old value and discard the new value
	Set(string) string
}

// OneTimeFactory is a function to handle the
// one time creation
type OneTimeFactory func(string, time.Duration) OneTimeSet

var (
	otFactory OneTimeFactory
)

// NewOneTimeSetter return a new one time setter
func NewOneTimeSetter(key string, d time.Duration) OneTimeSet {
	assert.NotNil(otFactory)
	return otFactory(key, d)
}
