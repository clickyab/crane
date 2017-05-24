package store

import (
	"time"

	"clickyab.com/exchange/services/assert"
)

// Interface is the blocking store interface
type Interface interface {
	// Push data in the store
	Push(string, string, time.Duration)
	// Pop and remove data from store, its blocking pop
	Pop(string, time.Duration) (string, bool)
}

// Factory is a function to handle the new store.Interface
type Factory func() Interface

var (
	factory Factory
)

// Register is a function to register a new factory
func Register(s Factory) {
	factory = s
}

// GetSyncStore return an in cluster sync
func GetSyncStore() Interface {
	assert.NotNil(factory, "[BUG] cluster factory is not set")
	return factory()
}
