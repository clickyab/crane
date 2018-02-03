package mock

import (
	"testing"

	"time"

	"github.com/clickyab/services/kv"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDlockSpec(t *testing.T) {
	kv.Register(nil, nil, NewMockDistributedLocker, nil, nil, nil, nil, nil)
	Convey("Test dlock locally", t, func() {

		Convey("test if the ttl is passed", func() {
			lock1 := kv.NewDistributedLock("TEST", 100*time.Millisecond)
			lock2 := kv.NewDistributedLock("TEST", 100*time.Millisecond)

			So(lock1.Resource(), ShouldEqual, "TEST")
			So(lock2.Resource(), ShouldEqual, "TEST")

			So(lock1.TTL(), ShouldEqual, 100*time.Millisecond)
			So(lock2.TTL(), ShouldEqual, 100*time.Millisecond)

			t := time.Now()
			lock1.Lock()
			lock2.Lock()
			// This should be called after a secound
			So(time.Since(t), ShouldBeGreaterThan, 100*time.Millisecond)
		})
		Convey("test if the time has not passed", func() {
			lock1 := kv.NewDistributedLock("TEST", time.Minute)
			lock2 := kv.NewDistributedLock("TEST", time.Minute)

			t := time.Now()
			lock1.Lock()
			lock1.Unlock()
			lock2.Lock()
			// This should be called after a secound
			So(time.Since(t), ShouldBeLessThan, time.Minute)
		})
	})
}
