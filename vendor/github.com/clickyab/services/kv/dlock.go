package kv

import (
	"sync"
	"time"

	"github.com/clickyab/services/assert"
)

// DistributedLock is the distributed locker interface.
type DistributedLock interface {
	sync.Locker
	// Resource return the resource to lock
	Resource() string
	// TTL return the ttl to wait for this lock
	TTL() time.Duration
}

// DistributedLockFactory is a factory function for lock
type DistributedLockFactory func(string, time.Duration) DistributedLock

var dlockFactory DistributedLockFactory

// NewDistributedLock return a new lock on resource with ttl
func NewDistributedLock(resource string, ttl time.Duration) DistributedLock {
	regLock.RLock()
	defer regLock.RUnlock()

	assert.NotNil(dlockFactory, "[BUG] no factory set for DistributedLock")
	return dlockFactory(resource, ttl)
}
