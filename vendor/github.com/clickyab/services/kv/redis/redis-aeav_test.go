package redis

import (
	"testing"
	"time"

	"github.com/clickyab/services/aredis"
	"github.com/clickyab/services/config"

	"github.com/clickyab/services/initializer"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAeavSpec(t *testing.T) {
	config.Initialize("test", "test", "test")
	defer initializer.Initialize()()

	Convey("Test keyvalue store for redis", t, func() {
		store := newRedisAEAVStore("atest_key", time.Second*100)
		So(store.Key(), ShouldEqual, "atest_key")
		Convey("check set and get", func() {
			store.IncSubKey("test", 1)
			So(store.SubKey("test"), ShouldEqual, 1)
			store.IncSubKey("another", 2)
			So(store.SubKey("another"), ShouldEqual, 2)
			So(store.AllKeys(), ShouldResemble, map[string]int64{"test": 1, "another": 2})
		})

		another := newRedisAEAVStore("atest_key", time.Second*100)
		So(another.SubKey("another"), ShouldEqual, 2)
		So(another.AllKeys(), ShouldResemble, map[string]int64{"test": 1, "another": 2})

		yet := newRedisAEAVStore("anot_exist", time.Second*100)
		So(yet.SubKey("anot_exist"), ShouldEqual, 0)

	})
	aredis.Client.FlushAll()

}
