package redis

import (
	"services/initializer"
	"services/redis"
	"testing"
	"time"

	config2 "services/config"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStoreRedis(t *testing.T) {
	config2.Initialize("test", "test", "test")
	defer initializer.Initialize()()
	Convey("test redis statistic store", t, func() {
		store := factory("test_static_key", 1*time.Hour)
		So(store.Key(), ShouldEqual, "test_static_key")
		Convey("inc and decr and touch sunbkey", func() {
			m, _ := store.IncSubKey("not_exists", 10)
			So(m, ShouldEqual, 10)
			n, _ := store.DecSubKey("not_exists", 5)
			So(n, ShouldEqual, 5)
			val, _ := store.Touch("not_exists")
			So(val, ShouldEqual, 5)
		})
		Convey("test get all key", func() {
			val, _ := store.GetAll()
			So(val["not_exists"], ShouldEqual, 5)
		})
	})
	aredis.Client.Del("test_static_key")
}
