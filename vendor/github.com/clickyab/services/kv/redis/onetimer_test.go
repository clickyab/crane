package redis

import (
	"testing"

	"time"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOneTimer(t *testing.T) {
	config.Initialize("test", "test", "test")
	defer initializer.Initialize()()

	Convey("test Key func", t, func() {
		key := "key1"
		// flush redis if u want to change value in 2 minutes
		store := newOneTimer(key, time.Minute*2)
		So(store.Key(), ShouldResemble, key)
	})

	Convey("test Set func", t, func() {
		key := "key1"
		value1 := "v1"
		value2 := "v2"

		store := newOneTimer(key, time.Minute*2)
		So(store.Key(), ShouldResemble, key)

		value := store.Set(value1)
		So(value, ShouldResemble, value1)
		value = store.Set(value2)
		So(value, ShouldResemble, value1)
	})
}
