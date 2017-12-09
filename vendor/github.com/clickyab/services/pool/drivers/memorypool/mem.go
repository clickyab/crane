package memorypool

import (
	"sync"

	"time"

	"fmt"

	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/pool"
)

type memPool struct {
	data map[string]kv.Serializable

	lock sync.RWMutex
}

func (mp *memPool) All() map[string]kv.Serializable {
	mp.lock.RLock()
	defer mp.lock.RUnlock()

	return mp.data
}

func (mp *memPool) Store(d map[string]kv.Serializable, _ time.Duration) error {
	mp.lock.Lock()
	defer mp.lock.Unlock()

	mp.data = d
	return nil
}

func (mp *memPool) Fetch(k string, data kv.Serializable) (kv.Serializable, error) {
	mp.lock.RLock()
	defer mp.lock.RUnlock()

	d, ok := mp.data[k]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	data = d
	return data, nil
}

// NewMemoryPool return an in memory pool
func NewMemoryPool() pool.Driver {
	return &memPool{}
}
