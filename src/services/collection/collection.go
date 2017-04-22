package collection

import "services/assert"

// Store is a collection store ready for to many write (big data i think)
type Store interface {
	// Name is the name of the store
	Name() string
	// Save try to save a collection in system
	Save(interface{}) error
}

// StoreFactory is
type StoreFactory func(string) Store

var factory StoreFactory

// Register is the base for registering a factory for this
func Register(f StoreFactory) {
	factory = f
}

// NewCollectionStore get the collection by its name
func NewCollectionStore(name string) Store {
	assert.NotNil(factory, "[BUG] no factory is registered")
	return factory(name)
}
