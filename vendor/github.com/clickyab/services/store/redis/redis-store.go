package redis

import (
	"time"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/redis"
	"github.com/clickyab/services/store"
)

type storeRedis struct {
}

// Push data in the store
func (storeRedis) Push(key string, value string, t time.Duration) {
	err := aredis.Client.LPush(key, value).Err()
	if err == nil {
		err = aredis.Client.Expire(key, t).Err()
	}

	assert.Nil(err)
}

// Pop and remove data from store, its blocking pop
func (storeRedis) Pop(key string, t time.Duration) (string, bool) {
	res := aredis.Client.BRPop(t, key)

	v := res.Val()
	if len(v) == 0 {
		return "", false
	}

	if len(v) == 2 && v[0] == key {
		return v[1], true
	}

	return "", false
}

func newRedisStore() store.Interface {
	return storeRedis{}
}

func init() {
	store.Register(newRedisStore)
}
