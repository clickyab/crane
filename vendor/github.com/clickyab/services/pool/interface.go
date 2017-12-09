package pool

import (
	"context"
	"time"

	"github.com/clickyab/services/kv"
)

// Loader is a function to handle the loading from any slow source
type Loader func(context.Context) (map[string]kv.Serializable, error)

// Driver is the storage driver of the pool
type Driver interface {
	// Store is the function to store data
	Store(map[string]kv.Serializable, time.Duration) error
	// Fetch try to fetch a single data
	Fetch(string, kv.Serializable) (kv.Serializable, error)

	All() map[string]kv.Serializable
}

// Interface is the pool interface
type Interface interface {
	// Get a single walue
	Get(string, kv.Serializable) (kv.Serializable, error)
	// All return all data if driver support it
	All() map[string]kv.Serializable
	// Start start the loading process
	Start(context.Context) context.Context
	// Notify is a hack. so we can wait for the first time.
	// TODO : watch it since after first time it may block the caller
	Notify() <-chan time.Time
}
