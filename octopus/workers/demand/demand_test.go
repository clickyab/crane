package demand

import (
	"testing"
	"time"

	"clickyab.com/exchange/octopus/exchange"
	"clickyab.com/exchange/octopus/exchange/materialize"
	"clickyab.com/exchange/octopus/workers/internal/datamodels"
	"clickyab.com/exchange/octopus/workers/mocks"
	"clickyab.com/exchange/services/assert"
	"clickyab.com/exchange/services/broker"

	"context"

	"clickyab.com/exchange/services/config"
	"clickyab.com/exchange/services/random"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	t1, _ = time.Parse("2006-01-02T15:04:05.000Z", "2017-03-21T00:01:00.000Z")
	//t2, _ = time.Parse("2006-01-02T15:04:05.000Z", "2017-03-21T01:01:00.000Z")
	//t3, _ = time.Parse("2006-01-02T15:04:05.000Z", "2017-03-21T02:01:00.000Z")
)

func newDemand(name string, rate int, handicap int64) exchange.Demand {
	return &mocks.Demand{
		DName:     name,
		DCallRate: rate,
		DHandicap: handicap,
	}
}

func newImpression(t time.Time, slotCount int, source, sup string) exchange.Impression {
	a := make([]mocks.Slot, 0)
	for i := 1; i <= slotCount; i++ {
		a = append(a, mocks.Slot{
			SWidth:   300,
			SHeight:  250,
			STRackID: <-random.ID,
		})
	}
	return mocks.Impression{
		ITime: t,
		ISource: mocks.Publisher{
			PName: source,
			PSupplier: mocks.Supplier{
				SName: sup,
			},
		},

		ISlots: a,
	}
}

func newAds(slots []exchange.Slot, demand exchange.Demand) map[string]exchange.Advertise {
	a := make(map[string]exchange.Advertise, 0)
	for i := range slots {
		a[slots[i].TrackID()] = &mocks.Ads{
			AHeight: slots[i].Height(),
			AWidth:  slots[i].Width(),
			AMaxCPM: 340,
			ADemand: demand,
		}
	}
	return a
}

func demToDelivery(i exchange.Impression, dem exchange.Demand, ads map[string]exchange.Advertise) broker.Delivery {
	job := materialize.DemandJob(i, dem, ads)
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

func TestDemand(t *testing.T) {
	config.Initialize("test", "test", "test")
	a := &agg{c: make(chan datamodels.TableModel, 2)}
	datamodels.RegisterAggregator(a)
	base := context.Background()
	Convey("demand json job", t, func() {
		d := newDemand("test_demand", 100, 50)
		imp := newImpression(t1, 2, "test_source", "test_supplier")
		//slots:=newSlots(2)
		ads := newAds(imp.Slots(), d)
		ctx, cnl := context.WithCancel(base)
		defer cnl()
		dem := consumer{ctx: ctx}
		delivery := dem.Consume()
		data := demToDelivery(imp, d, ads)
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
		So(t.Demand, ShouldEqual, "test_demand")
		So(t.Time, ShouldEqual, 1)
		So(t.DemandRequest, ShouldEqual, 1)
		So(t.ImpressionSlots, ShouldEqual, 0)
		So(t.ImpressionBid, ShouldEqual, 680)
		So(t.Source, ShouldEqual, "test_source")
		So(t.Supplier, ShouldEqual, "test_supplier")
	})

}
