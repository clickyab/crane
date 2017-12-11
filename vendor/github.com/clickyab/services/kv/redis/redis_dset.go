package redis

import (
	"time"

	"sync"

	"github.com/clickyab/services/aredis"
	"github.com/clickyab/services/kv"
)

type dsetRedis struct {
	once    sync.Once
	key     string
	members []string
	add     []string
	locker  sync.Mutex
}

// Members return all members of this set
func (k *dsetRedis) Members() []string {
	k.locker.Lock()
	defer k.locker.Unlock()
	k.once.Do(func() {
		c := aredis.Client.SMembers(k.key)
		if c.Err() != nil {
			return
		}
		res, err := c.Result()
		if err != nil {
			return
		}
		k.members = res
	})
	return append(k.members, k.add...)
}

// Add new item to set
func (k *dsetRedis) Add(s ...string) {
	k.locker.Lock()
	defer k.locker.Unlock()
	k.add = append(k.add, s...)
}

// Key return the master key
func (k *dsetRedis) Key() string {
	return k.key
}

// Save the set with lifetime
func (k *dsetRedis) Save(t time.Duration) error {
	k.locker.Lock()
	defer k.locker.Unlock()
	dt := make([]interface{}, 0)
	for _, s := range k.add {
		dt = append(dt, s)
	}
	res := aredis.Client.SAdd(k.Key(), dt...)
	if res.Err() != nil {
		return res.Err()
	}
	k.add = nil
	aredis.Client.Expire(k.key, t)
	return res.Err()
}

// newRedisDsetStore return a redis store for eav
func newRedisDsetStore(key string) kv.DistributedSet {
	return &dsetRedis{
		key: key,
	}
}
