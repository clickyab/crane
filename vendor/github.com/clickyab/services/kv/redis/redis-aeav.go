package redis

import (
	"sync"
	"time"

	"strconv"

	"github.com/clickyab/services/aredis"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/kv"
)

type atomicKiwiRedis struct {
	key string
	v   map[string]int64
	sync.Mutex
	duration time.Duration
}

func (kr *atomicKiwiRedis) TTL() time.Duration {
	d := aredis.Client.TTL(kr.key)
	r, _ := d.Result()

	return r
}

func (kr *atomicKiwiRedis) Drop(s ...string) error {
	kr.Lock()
	defer kr.Unlock()

	if len(s) == 0 {
		kr.v = make(map[string]int64)
		d := aredis.Client.Del(kr.key)
		return d.Err()
	}

	for i := range s {
		delete(kr.v, s[i])
	}
	// Ignore the error
	d := aredis.Client.HDel(kr.key, s...)
	return d.Err()
}

// Key return the parent key
func (kr *atomicKiwiRedis) Key() string {
	return kr.key
}

// IncSubKey for increasing sub key
func (kr *atomicKiwiRedis) IncSubKey(key string, value int64) int64 {
	res := aredis.Client.HIncrBy(kr.key, key, value)
	if res.Err() != nil {
		kr.v[key] = value
		return kr.v[key]
	}
	r, err := res.Result()
	if err != nil {
		kr.v[key] = value
		return r
	}
	kr.v[key] = r
	aredis.Client.Expire(key, kr.duration)
	return r
}

// IncSubKey for decreasing sub key
func (kr *atomicKiwiRedis) DecSubKey(key string, value int64) int64 {
	return kr.IncSubKey(key, -value)
}

// SubKey return a key
func (kr *atomicKiwiRedis) SubKey(key string) int64 {
	kr.Lock()
	defer kr.Unlock()

	if v, ok := kr.v[key]; ok {
		aredis.Client.Expire(key, kr.duration)
		return v
	}
	res := aredis.Client.HIncrBy(kr.key, key, 0)
	if res.Err() != nil {
		return 0
	}

	r, err := res.Result()
	if err != nil {
		return 0
	}

	return r
}

// AllKeys from the store
func (kr *atomicKiwiRedis) AllKeys() map[string]int64 {
	kr.v = map[string]int64{}
	res := aredis.Client.HGetAll(kr.key)

	if res.Err() != nil {
		return kr.v
	}

	r, err := res.Result()
	if err != nil {
		return kr.v
	}
	f := make(map[string]int64)
	for k, v := range r {
		tv, e := strconv.ParseInt(v, 10, 64)
		assert.Nil(e)
		f[k] = tv
	}
	kr.v = f
	return kr.v
}

// NewRedisAEAVStore return a redis store for eav
func newRedisAEAVStore(key string, dur time.Duration) kv.AKiwi {
	return &atomicKiwiRedis{
		key:      key,
		v:        make(map[string]int64),
		Mutex:    sync.Mutex{},
		duration: dur,
	}
}
