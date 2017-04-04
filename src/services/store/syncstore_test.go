package store

import (
	"testing"

	"time"

	"services/store/mock"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSyncStore(t *testing.T) {
	Register(func() Interface {
		return mock.NewMockChannelStore()
	})
	Convey("Test simple push/pop", t, func() {
		tmp := GetSyncStore()
		tmp.Push("test", "value", time.Minute)
		v, ok := tmp.Pop("test", time.Millisecond)
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "value")

		v, ok = tmp.Pop("notvalidkey", time.Millisecond)
		So(ok, ShouldBeFalse)
		So(v, ShouldBeEmpty)
	})
}
