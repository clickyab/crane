package mock

import (
	"testing"

	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSyncStore(t *testing.T) {

	Convey("Test simple push/pop", t, func() {
		tmp := NewMockChannelStore()
		tmp.Push("test", "value", time.Minute)
		v, ok := tmp.Pop("test", time.Millisecond)
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "value")

		v, ok = tmp.Pop("notvalidkey", time.Millisecond)
		So(ok, ShouldBeFalse)
		So(v, ShouldBeEmpty)
	})
}
