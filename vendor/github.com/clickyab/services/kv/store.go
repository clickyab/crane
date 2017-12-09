package kv

import (
	"time"

	"github.com/clickyab/services/assert"
)

// Store is the blocking store interface
type Store interface {
	// Push data in the store
	Push(string, string, time.Duration)
	// Pop and remove data from store, its blocking pop
	Pop(string, time.Duration) (string, bool)
}

// StoreFactory is a function to handle the new store.Interface
type StoreFactory func() Store

var storeFactory StoreFactory

// GetSyncStore return an in cluster sync
func GetSyncStore() Store {
	regLock.RLock()
	defer regLock.RUnlock()

	assert.NotNil(storeFactory, "[BUG] cluster factory is not set")
	return storeFactory()
}
