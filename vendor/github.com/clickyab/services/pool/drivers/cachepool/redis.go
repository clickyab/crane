package cachepool

import (
	"time"

	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/pool"
)

type redisPool struct {
	prefix string
}

func (rp *redisPool) All() map[string]kv.Serializable {
	panic("not supported on this driver")
}

func (rp *redisPool) Store(d map[string]kv.Serializable, t time.Duration) error {
	for i := range d {
		err := kv.Do(rp.prefix+i, d[i], t, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rp *redisPool) Fetch(k string, data kv.Serializable) (kv.Serializable, error) {
	k = rp.prefix + k
	if err := kv.Hit(k, data); err != nil {
		return nil, err
	}
	return data, nil
}

// NewCachePool return a new cache pool
func NewCachePool(prefix string) pool.Driver {
	return &redisPool{
		prefix: prefix,
	}
}
