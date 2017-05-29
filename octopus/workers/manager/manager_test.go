package manager

import (
	"testing"

	"clickyab.com/exchange/octopus/workers/internal/datamodels"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAggregateFunc(t *testing.T) {
	Convey("the demand test with the impression job", t, func() {
		a := datamodels.TableModel{
			Time:               1,
			Demand:             "Demand",
			Supplier:           "Supplier",
			Source:             "Source",
			DeliverCount:       1,
			DeliverBid:         100,
			WinCount:           2,
			WinBid:             200,
			RequestOutCount:    3,
			RequestInCount:     4,
			ImpressionInCount:  6,
			ImpressionOutCount: 2,
		}

		b := aggregate(nil, a)
		So(*b, ShouldResemble, a)
		c := aggregate(b, a)

		So(c.Time, ShouldEqual, a.Time)
		So(c.Demand, ShouldEqual, a.Demand)
		So(c.Supplier, ShouldEqual, a.Supplier)
		So(c.Source, ShouldEqual, a.Source)
		So(c.DeliverCount, ShouldEqual, 2)
		So(c.DeliverBid, ShouldEqual, 200)
		So(c.WinCount, ShouldEqual, 4)
		So(c.WinBid, ShouldEqual, 400)
		So(c.ImpressionInCount, ShouldEqual, 12)
		So(c.ImpressionOutCount, ShouldEqual, 4)
		So(c.RequestOutCount, ShouldEqual, 6)
		So(c.RequestInCount, ShouldEqual, 8)

		a.Time = 3
		So(func() { aggregate(c, a) }, ShouldPanic)
	})
}
