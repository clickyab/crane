package manager

import (
	"testing"

	"clickyab.com/exchange/octopus/workers/internal/datamodels"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAggregateFunc(t *testing.T) {
	Convey("the demand test with the impression job", t, func() {
		a := datamodels.TableModel{
			Time:              1,
			Demand:            "Demand",
			Supplier:          "Supplier",
			Source:            "Source",
			ShowCount:         1,
			ShowBid:           100,
			WinCount:          2,
			WinBid:            200,
			ImpressionSlots:   3,
			ImpressionBid:     300,
			ImpressionRequest: 4,
		}

		b := aggregate(nil, a)
		So(*b, ShouldResemble, a)
		c := aggregate(b, a)

		So(c.Time, ShouldEqual, a.Time)
		So(c.Demand, ShouldEqual, a.Demand)
		So(c.Supplier, ShouldEqual, a.Supplier)
		So(c.Source, ShouldEqual, a.Source)
		So(c.ShowCount, ShouldEqual, 2)
		So(c.ShowBid, ShouldEqual, 200)
		So(c.WinCount, ShouldEqual, 4)
		So(c.WinBid, ShouldEqual, 400)
		So(c.ImpressionSlots, ShouldEqual, 6)
		So(c.ImpressionBid, ShouldEqual, 600)

		a.Time = 3
		So(func() { aggregate(c, a) }, ShouldPanic)
	})
}
