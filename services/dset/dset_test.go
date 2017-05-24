package dset_test

import (
	"testing"
	"time"

	. "clickyab.com/exchange/services/dset"
	"clickyab.com/exchange/services/dset/mock"

	"github.com/smartystreets/goconvey/convey"
)

func TestManager(t *testing.T) {
	Register(mock.NewMockDsetStore)
	convey.Convey("test set", t, func() {
		d := NewDistributedSet("test_key")
		convey.So(d.Key(), convey.ShouldEqual, "test_key")
		d.Add("1", "2")
		convey.So(len(d.Members()), convey.ShouldEqual, 2)
		d.Save(2 * time.Second)
		convey.So(d.Members(), convey.ShouldResemble, []string{"1", "2"})
		e := NewDistributedSet("test_key")
		convey.So(d.Members(), convey.ShouldResemble, e.Members())
		time.Sleep(3 * time.Second)
		f := NewDistributedSet("test_key")
		convey.So(len(f.Members()), convey.ShouldEqual, 0)
	})
}
