package broker

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type dum struct {
}

func (dum) Publish(Job) {
	panic("HI THIS IS ME")
}

func TestSpec(t *testing.T) {
	Convey("Test register publisher ", t, func() {
		So(func() {
			SetActiveBroker(dum{})
		}, ShouldNotPanic)
		So(func() {
			SetActiveBroker(dum{})
		}, ShouldPanic)

		So(func() { Publish(nil) }, ShouldPanicWith, "HI THIS IS ME")
	})
}
