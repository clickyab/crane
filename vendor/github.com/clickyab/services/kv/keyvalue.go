package kv

import (
	"time"

	"github.com/clickyab/services/assert"
)

// Kiwi is the key value storage system in a parent key
type Kiwi interface {
	// Key return the parent key
	Key() string
	// SetSubKey for adding a sub key
	SetSubKey(key, value string) Kiwi
	// SubKey return a key
	SubKey(key string) string
	// AllKeys from the store
	AllKeys() map[string]string
	// Save the entire keys (mostly first time)
	Save(time.Duration) error
	// Drop the entire eav store
	Drop() error
	// TTL return the time to expire this
	TTL() time.Duration
}

// KiwiFactory is a function to create store
type KiwiFactory func(string) Kiwi

var (
	kiwiFactory KiwiFactory
)

// NewEavStore return a new eav store
func NewEavStore(key string) Kiwi {
	regLock.RLock()
	defer regLock.RUnlock()

	assert.NotNil(kiwiFactory, "[BUG] factory is not registered")
	return kiwiFactory(key)
}
