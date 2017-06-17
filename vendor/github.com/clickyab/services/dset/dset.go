package dset

import (
	"time"

	"github.com/clickyab/services/assert"
)

// DistributedSet is a set distributed
type DistributedSet interface {
	// Members return all members of this set
	Members() []string
	// Add new item to set
	Add(...string)
	// Key return the master key
	Key() string
	// Save the set with lifetime
	Save(time.Duration) error
}

// DistributedSetFactory is a function to create store
type DistributedSetFactory func(string) DistributedSet

var (
	factory DistributedSetFactory
)

// Register is a function to register store factory
func Register(s DistributedSetFactory) {
	factory = s
}

// NewDistributedSet is the distributed set
func NewDistributedSet(key string) DistributedSet {
	assert.NotNil(factory, "[BUG] factory is not registered")
	return factory(key)
}
