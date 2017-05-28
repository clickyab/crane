package show

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"time"

	"context"

	"clickyab.com/exchange/octopus/workers/internal/datamodels"
	"clickyab.com/exchange/octopus/workers/mocks"
	"clickyab.com/exchange/services/broker"
	"clickyab.com/exchange/services/config"
)

var (
	raw = `{"ad_id":"64675","demand_name":"clickyab-demo","price":5050,"publisher":"entekhab.ir","slot_id":"e14bd6611055d69fb8883a016474e99787175f37","supplier":"clickyab","time":"1495869393","track_id":"e14bd6611055d69fb8883a016474e99787175f37"}`
)

type agg struct {
	c chan datamodels.TableModel
}

func (a *agg) Channel() chan<- datamodels.TableModel {
	return a.c
}

func winToDelivery() broker.Delivery {
	return mocks.JsonDelivery{Data: []byte(raw)}
}

func TestImpression(t *testing.T) {
	config.Initialize("test", "test", "test")
	a := &agg{c: make(chan datamodels.TableModel, 2)}
	datamodels.RegisterAggregator(a)
	base := context.Background()
	Convey("the demand test with the show job", t, func() {
		ctx, cl := context.WithCancel(base)
		defer cl()
		dem := consumer{ctx: ctx}
		delivery := dem.Consume()
		data := winToDelivery()
		select {
		case delivery <- data:
			So(true, ShouldBeTrue)
		case <-time.After(time.Second):
			So(true, ShouldBeFalse)
		}
		var t datamodels.TableModel
		select {
		case t = <-a.c:
			So(true, ShouldBeTrue)
		case <-time.After(time.Second):
			So(true, ShouldBeFalse)
		}
		select {
		case <-a.c:
			So(true, ShouldBeFalse)
		case <-time.After(time.Second):
			So(true, ShouldBeTrue)
		}

		So(t.Supplier, ShouldEqual, "clickyab")
		So(t.Source, ShouldEqual, "entekhab.ir")
		So(t.Demand, ShouldEqual, "clickyab-demo")
		So(t.ShowCount, ShouldEqual, 1)
		So(t.ShowBid, ShouldEqual, 5050)
	})
}
