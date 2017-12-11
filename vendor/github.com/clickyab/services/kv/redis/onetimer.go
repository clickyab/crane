package redis

import (
	"time"

	"github.com/clickyab/services/aredis"
	"github.com/clickyab/services/kv"
)

type oneTimer struct {
	key string
	d   time.Duration
}

func (ot *oneTimer) Key() string {
	return ot.key
}

func (ot *oneTimer) Set(s string) string {
	b := aredis.Client.SetNX(ot.key, s, ot.d)
	if b.Err() != nil {
		return s
	}
	if b.Val() { // means it set the value so its for the first time
		return s
	}

	v := aredis.Client.Get(ot.key)
	if v.Err() != nil {
		return s
	}
	_ = aredis.Client.Expire(ot.key, ot.d) // ignore error
	return v.Val()
}

func newOneTimer(key string, d time.Duration) kv.OneTimeSet {
	return &oneTimer{key, d}
}
