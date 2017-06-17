package statistic

import (
	"time"

	"github.com/clickyab/services/assert"
)

// Interface is the base interface to handle the statistics key
type Interface interface {
	// Key return the master key. the key is important to group the subkeys
	Key() string
	// IncSubKey increase the value for a sub key and return the value
	IncSubKey(string, int64) (int64, error)
	// IncSubKey decrease the value for a sub key and return the value
	DecSubKey(string, int64) (int64, error)
	// Touch return the current value of the sub key
	Touch(string) (int64, error)
	// GetAll return all keys in single call
	GetAll() (map[string]int64, error)
}

// Factory is a function to handle the new store.Interface
type Factory func(string, time.Duration) Interface

var (
	factory Factory
)

// Register is a function to register a new factory
func Register(s Factory) {
	factory = s
}

// GetStatisticStore return an in cluster sync
func GetStatisticStore(key string, expire time.Duration) Interface {
	assert.NotNil(factory, "[BUG] cluster factory is not set")
	return factory(key, expire)
}
