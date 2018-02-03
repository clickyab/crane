package mock

import (
	"testing"

	"time"

	"reflect"

	"github.com/clickyab/services/kv"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {
	kv.Register(NewMockStore, nil, nil, nil, nil, nil, nil, nil)

	Convey("Test keyvalue store", t, func() {
		store := kv.NewEavStore("test_key")
		So(store.Key(), ShouldEqual, "test_key")
		Convey("check set and get", func() {
			store.SetSubKey("test", "test_val")
			So(store.SubKey("test"), ShouldEqual, "test_val")
			store.SetSubKey("another", "2")
			So(store.SubKey("another"), ShouldEqual, "2")
			So(store.Save(time.Hour), ShouldBeNil)
			So(reflect.DeepEqual(store.AllKeys(), map[string]string{"test": "test_val", "another": "2"}), ShouldBeTrue)
		})
	})
}
