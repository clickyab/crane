package kv

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
	dsetFactory DistributedSetFactory
)

// NewDistributedSet is the distributed set
func NewDistributedSet(key string) DistributedSet {
	regLock.RLock()
	defer regLock.RUnlock()

	assert.NotNil(dsetFactory, "[BUG] factory is not registered")
	return dsetFactory(key)
}
