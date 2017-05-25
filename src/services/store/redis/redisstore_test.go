package redis

import (
	"services/config"
	"services/initializer"
	"services/redis"
	"services/store"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSyncStore(t *testing.T) {
	config.Initialize("test", "test", "test")
	defer initializer.Initialize()()
	Convey("Test simple push/pop", t, func() {
		tmp := store.GetSyncStore()
		tmp.Push("test_random_value", "value", time.Minute)
		v, ok := tmp.Pop("test_random_value", time.Millisecond)
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "value")

		v, ok = tmp.Pop("notvalidkey", time.Millisecond)
		So(ok, ShouldBeFalse)
		So(v, ShouldBeEmpty)
		aredis.Client.Del("test_random_value")
	})
}
