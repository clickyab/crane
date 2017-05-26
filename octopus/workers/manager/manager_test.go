package manager

import (
	"testing"

	"clickyab.com/exchange/octopus/workers/internal/datamodels"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAggregateFunc(t *testing.T) {
	Convey("the demand test with the impression job", t, func() {
		a := datamodels.TableModel{
			Time:          1,
			Demand:        "Demand",
			Supplier:      "Supplier",
			Source:        "Source",
			Show:          1,
			ShowBid:       100,
			Win:           2,
			WinnerBid:     200,
			Impression:    3,
			ImpressionBid: 300,
			Request:       4,
		}

		b := aggregate(nil, a)
		So(*b, ShouldResemble, a)
		c := aggregate(b, a)

		So(c.Time, ShouldEqual, a.Time)
		So(c.Demand, ShouldEqual, a.Demand)
		So(c.Supplier, ShouldEqual, a.Supplier)
		So(c.Source, ShouldEqual, a.Source)
		So(c.Show, ShouldEqual, 2)
		So(c.ShowBid, ShouldEqual, 200)
		So(c.Win, ShouldEqual, 4)
		So(c.WinnerBid, ShouldEqual, 400)
		So(c.Impression, ShouldEqual, 6)
		So(c.ImpressionBid, ShouldEqual, 600)

		a.Time = 3
		So(func() { aggregate(c, a) }, ShouldPanic)
	})
}
