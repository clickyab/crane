package store

import (
	"testing"

	"time"

	"clickyab.com/exchange/services/store/mock"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSyncStore(t *testing.T) {
	Register(func() Interface {
		return mock.NewMockChannelStore()
	})
	Convey("Test simple push/pop", t, func() {
		tmp := GetSyncStore()
		tmp.Push("test_push_value", "value", time.Minute)
		v, ok := tmp.Pop("test_push_value", time.Millisecond)
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "value")

		v, ok = tmp.Pop("notvalidkey_test_push_value", time.Millisecond)
		So(ok, ShouldBeFalse)
		So(v, ShouldBeEmpty)
	})
}
