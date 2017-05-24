package dlock

import (
	"sync"
	"time"

	"clickyab.com/exchange/services/assert"
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

var factory DistributedLockFactory

// Register is a function to register store factory
func Register(s DistributedLockFactory) {
	factory = s
}

// NewDistributedLock return a new lock on resource with ttl
func NewDistributedLock(resource string, ttl time.Duration) DistributedLock {
	assert.NotNil(factory, "[BUG] no factory set for DistributedLock")
	return factory(resource, ttl)
}
