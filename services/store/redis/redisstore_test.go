package redis

import (
	"testing"
	"time"

	"clickyab.com/exchange/services/config"
	"clickyab.com/exchange/services/initializer"
	"clickyab.com/exchange/services/redis"
	"clickyab.com/exchange/services/store"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSyncStore(t *testing.T) {
	config.Initialize("test", "test", "test")
	defer initializer.Initialize()()
	Convey("Test simple push/pop", t, func() {
		tmp := store.GetSyncStore()
		tmp.Push("test_redis_store", "value", time.Minute)
		v, ok := tmp.Pop("test_redis_store", time.Millisecond)
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "value")

		v, ok = tmp.Pop("notvalidkey_test_redis_store", time.Millisecond)
		So(ok, ShouldBeFalse)
		So(v, ShouldBeEmpty)
		aredis.Client.Del("test_redis_store")
	})
}
