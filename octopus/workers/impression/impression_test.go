package impression

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"time"

	"context"

	"clickyab.com/exchange/octopus/exchange"
	"clickyab.com/exchange/octopus/exchange/materialize"
	"clickyab.com/exchange/octopus/workers/internal/datamodels"
	"clickyab.com/exchange/octopus/workers/mocks"
	"clickyab.com/exchange/services/assert"
	"clickyab.com/exchange/services/broker"
	"clickyab.com/exchange/services/config"
)

var (
	t1, _ = time.Parse("2006-01-02T15:04:05.000Z", "2017-03-21T00:01:00.000Z")
	//t2, _ = time.Parse("2006-01-02T15:04:05.000Z", "2017-03-21T01:01:00.000Z")
	//t3, _ = time.Parse("2006-01-02T15:04:05.000Z", "2017-03-21T02:01:00.000Z")
)

func newImpression(t time.Time, slotCount int, source, sup string) exchange.Impression {
	return mocks.Impression{
		ITime: t,
		ISource: mocks.Publisher{
			PName: source,
			PSupplier: mocks.Supplier{
				SName: sup,
			},
		},
		ISlots: make([]mocks.Slot, slotCount),
	}
}

func impToDelivery(in exchange.Impression) broker.Delivery {
	job := materialize.ImpressionJob(in)
	d, err := job.Encode()
	assert.Nil(err)
	return mocks.JsonDelivery{Data: d}
}

type agg struct {
	c chan datamodels.TableModel
}

func (a *agg) Channel() chan<- datamodels.TableModel {
	return a.c
}

func TestImpression(t *testing.T) {
	config.Initialize("test", "test", "test")
	// make sure channel has space for more than 1 delivery
	a := &agg{c: make(chan datamodels.TableModel, 2)}
	datamodels.RegisterAggregator(a)
	base := context.Background()
	Convey("the demand test with the impression job", t, func() {
		imp := newImpression(t1, 10, "test_source", "test_sup")
		ctx, cl := context.WithCancel(base)
		defer cl()
		dem := consumer{ctx: ctx}

		delivery := dem.Consume()
		data := impToDelivery(imp)
		// make sure this is not blocker, and if the test fails then may it hangs for ever
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

		So(t.Supplier, ShouldEqual, "test_sup")
		So(t.Source, ShouldEqual, "test_source")
		So(t.Time, ShouldEqual, 1)
		So(t.Request, ShouldEqual, 1)
		So(t.Impression, ShouldEqual, 10)

		So(t.WinnerBid, ShouldBeZeroValue)
		So(t.ShowBid, ShouldBeZeroValue)
		So(t.Show, ShouldBeZeroValue)
		So(t.Demand, ShouldBeZeroValue)
	})
}
