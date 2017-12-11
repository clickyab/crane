package redis

import (
	"time"

	"github.com/clickyab/services/aredis"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/kv"
)

type psRedis struct {
}

// Push data in the store
func (psRedis) Push(key string, value string, t time.Duration) {
	err := aredis.Client.LPush(key, value).Err()
	if err == nil {
		err = aredis.Client.Expire(key, t).Err()
	}

	assert.Nil(err)
}

// Pop and remove data from store, its blocking pop
func (psRedis) Pop(key string, t time.Duration) (string, bool) {
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

func newRedisStore() kv.Store {
	return psRedis{}
}
