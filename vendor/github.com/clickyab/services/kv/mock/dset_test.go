package mock

import (
	"testing"
	"time"

	"github.com/clickyab/services/kv"
	"github.com/smartystreets/goconvey/convey"
)

func TestManager(t *testing.T) {
	kv.Register(nil, nil, nil, NewMockDsetStore, nil, nil, nil, nil)

	convey.Convey("test set", t, func() {
		d := kv.NewDistributedSet("test_key")
		convey.So(d.Key(), convey.ShouldEqual, "test_key")
		d.Add("1", "2")
		convey.So(len(d.Members()), convey.ShouldEqual, 2)
		d.Save(2 * time.Second)
		convey.So(d.Members(), convey.ShouldResemble, []string{"1", "2"})
		e := kv.NewDistributedSet("test_key")
		convey.So(d.Members(), convey.ShouldResemble, e.Members())
		time.Sleep(3 * time.Second)
		f := kv.NewDistributedSet("test_key")
		convey.So(len(f.Members()), convey.ShouldEqual, 0)
	})
}
