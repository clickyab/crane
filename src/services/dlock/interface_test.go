package dlock_test

import (
	. "services/dlock"
	"services/dlock/mock"
	"testing"

	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {
	Register(mock.NewMockDistributedLocker)
	Convey("Test dlock locally", t, func() {

		Convey("test if the ttl is passed", func() {
			lock1 := NewDistributedLock("TEST", 100*time.Millisecond)
			lock2 := NewDistributedLock("TEST", 100*time.Millisecond)

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
			lock1 := NewDistributedLock("TEST", time.Minute)
			lock2 := NewDistributedLock("TEST", time.Minute)

			t := time.Now()
			lock1.Lock()
			lock1.Unlock()
			lock2.Lock()
			// This should be called after a secound
			So(time.Since(t), ShouldBeLessThan, time.Minute)
		})
	})
}
