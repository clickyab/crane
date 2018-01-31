package mock

import (
	"time"

	"github.com/clickyab/services/kv"
)

// Kiwi is a mock kiwi
type Kiwi struct {
	MasterKey string
	Data      map[string]string
	Duration  time.Duration
}

var (
	pre   = make(map[string]map[string]string)
	store = make(map[string]*Kiwi)
)

// Drop the key
func (m *Kiwi) Drop() error {
	m.Data = make(map[string]string)
	pre[m.MasterKey] = nil

	return nil
}

// Key return the parent key
func (m *Kiwi) Key() string {
	return m.MasterKey
}

// SetSubKey for adding a sub key
func (m *Kiwi) SetSubKey(key, value string) kv.Kiwi {
	m.Data[key] = value
	return m
}

// SubKey return a key
func (m *Kiwi) SubKey(key string) string {
	return m.Data[key]
}

// AllKeys from the store
func (m *Kiwi) AllKeys() map[string]string {
	return m.Data
}

// Save the entire keys (mostly first time)
func (m *Kiwi) Save(t time.Duration) error {
	m.Duration = t
	return nil
}

// TTL return the expiration time of this
func (m *Kiwi) TTL() time.Duration {
	return m.Duration
}

// NewMockStore is the new mock store
func NewMockStore(key string) kv.Kiwi {
	if k, ok := store[key]; ok {
		return k
	}
	var (
		data map[string]string
		ok   bool
	)
	if data, ok = pre[key]; !ok {
		data = make(map[string]string)
	}
	m := &Kiwi{
		MasterKey: key,
		Data:      data,
	}

	store[key] = m
	return m
}

// SetMockData try to set mock data if needed
func SetMockData(key string, data map[string]string) {
	pre[key] = data
}

// GetMockStore is a function to get the mock store for testing
func GetMockStore() map[string]map[string]string {
	res := make(map[string]map[string]string)
	for i := range store {
		res[i] = store[i].Data
	}

	return res
}

// ResetEav the entire mock
func ResetEav() {
	store = make(map[string]*Kiwi)
}
