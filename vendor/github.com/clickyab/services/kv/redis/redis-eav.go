package redis

import (
	"sync"
	"time"

	"github.com/clickyab/services/aredis"
	"github.com/clickyab/services/kv"
)

type kiwiRedis struct {
	key  string
	v    map[string]string
	lock sync.Mutex
}

func (kr *kiwiRedis) TTL() time.Duration {
	d := aredis.Client.TTL(kr.key)
	r, _ := d.Result()

	return r
}

func (kr *kiwiRedis) Drop() error {
	kr.lock.Lock()
	defer kr.lock.Unlock()

	kr.v = make(map[string]string)
	d := aredis.Client.Del(kr.key)
	return d.Err()
}

// Key return the parent key
func (kr *kiwiRedis) Key() string {
	return kr.key
}

// SetSubKey for adding a sub key
func (kr *kiwiRedis) SetSubKey(key, value string) kv.Kiwi {
	kr.lock.Lock()
	defer kr.lock.Unlock()

	kr.v[key] = value

	return kr
}

// SubKey return a key
func (kr *kiwiRedis) SubKey(key string) string {
	kr.lock.Lock()
	defer kr.lock.Unlock()

	if v, ok := kr.v[key]; ok {
		return v
	}
	res := aredis.Client.HGet(kr.key, key)
	if res.Err() != nil {
		return ""
	}

	r, err := res.Result()
	if err != nil {
		return ""
	}

	return r
}

// AllKeys from the store
func (kr *kiwiRedis) AllKeys() map[string]string {
	kr.lock.Lock()
	defer kr.lock.Unlock()

	kr.v = map[string]string{}
	res := aredis.Client.HGetAll(kr.key)
	if res.Err() != nil {
		return kr.v
	}

	r, err := res.Result()
	if err != nil {
		return kr.v
	}

	kr.v = r
	return kr.v
}

// Save the entire keys (mostly first time)
func (kr *kiwiRedis) Save(t time.Duration) error {
	kr.lock.Lock()
	defer kr.lock.Unlock()

	tmp := make(map[string]interface{}, len(kr.v))
	for i := range kr.v {
		tmp[i] = kr.v[i]
	}

	res := aredis.Client.HMSet(kr.key, tmp)
	if res.Err() != nil {
		return res.Err()
	}

	b := aredis.Client.Expire(kr.key, t)
	return b.Err()
}

// newRedisEAVStore return a redis store for eav
func newRedisEAVStore(key string) kv.Kiwi {
	return &kiwiRedis{
		key:  key,
		v:    make(map[string]string),
		lock: sync.Mutex{},
	}
}
