package eav

import (
	"services/assert"
	"time"
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
}

// StoreFactory is a function to create store
type StoreFactory func(string) Kiwi

var (
	factory StoreFactory
)

// Register is a function to register store factory
func Register(s StoreFactory) {
	factory = s
}

// NewEavStore return a new eav store
func NewEavStore(key string) Kiwi {
	assert.NotNil(factory, "[BUG] factory is not registered")
	return factory(key)
}
