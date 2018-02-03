package redis

import (
	"testing"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"

	"time"

	"sort"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRedis(t *testing.T) {
	config.Initialize("test", "test", "test")
	defer initializer.Initialize()()
	Convey("Add function test", t, func() {
		ds := newRedisDsetStore("test")
		ds.Add("a1", "a2")

		So(ds.Members(), ShouldResemble, []string{"a1", "a2"})
	})

	Convey("save function test", t, func() {
		ds := newRedisDsetStore("test")
		ds.Add("a1", "a2")
		ds.Save(5 * time.Second)
		ds.Add("a3")

		expect := ds.Members()
		sort.Strings(expect)

		So(expect, ShouldResemble, []string{"a1", "a2", "a3"})
	})

	Convey("expiration test", t, func() {
		ds := newRedisDsetStore("test")
		ds.Add("a1", "a2")
		ds.Save(1 * time.Second)
		ds.Add("a3")

		expect := ds.Members()
		sort.Strings(expect)

		time.Sleep(2 * time.Second)
		ps := newRedisDsetStore("test")

		So(len(ps.Members()), ShouldEqual, 0)
	})

	Convey("expiration test 2", t, func() {
		ds := newRedisDsetStore("test")
		ds.Add("a1", "a2")
		ds.Save(1 * time.Second)

		expect := ds.Members()
		sort.Strings(expect)

		ps := newRedisDsetStore("test")

		So(len(ps.Members()), ShouldEqual, len(ds.Members()))
	})
}
