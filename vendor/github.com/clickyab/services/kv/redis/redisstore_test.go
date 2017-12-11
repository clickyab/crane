package redis

import (
	"testing"
	"time"

	"github.com/clickyab/services/aredis"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"

	"github.com/clickyab/services/kv"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSyncStore(t *testing.T) {
	config.Initialize("test", "test", "test")
	defer initializer.Initialize()()
	Convey("Test simple push/pop", t, func() {
		tmp := kv.GetSyncStore()
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
