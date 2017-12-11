package redis

import (
	"github.com/clickyab/services/aredis"
	"github.com/clickyab/services/kv"
)

type register struct {
}

func (register) Initialize() {
	kv.Register(newRedisEAVStore, newRedisStore, newRedisDistributedLock, newRedisDsetStore, newRedisAEAVStore, &cache{}, newRedisScanner, newOneTimer)
}

func init() {
	aredis.Register(register{})
}
