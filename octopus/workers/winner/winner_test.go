package winner

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"time"

	"context"

	"clickyab.com/exchange/octopus/workers/internal/datamodels"
	"clickyab.com/exchange/octopus/workers/mocks"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/config"
)

var (
	//winner_test
	//t1, _ = time.Parse("2006-01-02T15:04:05.000Z", "2017-03-21T00:01:00.000Z")
	raw = `{"advertise":{"demand":{"call_rate":1,"excluded_suppliers":[],"handicap":100,"name":"clickyab-demo","white_list_countries":[]},"height":250,"id":"64681","landing":"www.salamateaval.com","max_cpm":5050,"rate":null,"slot_track_id":"92005149213","track_id":"a5386297e19f1fa5cba3e6c745f88f7f0d09804a","url":"http://e.clickyab.com/show/sync/8bc27d90a0948342091dfa50cd5e253bcd62a6bc/35997/m184df33b0371e0cd392ee3307868fce792f52675c?parent=\u0026ref=\u0026s=3856691\u0026tid=53890855","width":300,"winner_cpm":5050},"impression":{"attributes":null,"category":[],"ip":"46.209.239.50","location":{"country":{"valid":true,"name":"Iran, Islamic Republic of","iso":"IR"},"lat_lon":{"valid":false,"lat":0,"lon":0},"province":{"valid":true,"name":"Tehran"}},"platform":2,"slots":[{"ad":{"demand":{"call_rate":1,"excluded_suppliers":[],"handicap":100,"name":"clickyab-demo","white_list_countries":[]},"height":250,"id":"64681","landing":"www.salamateaval.com","max_cpm":5050,"rate":null,"slot_track_id":"92005149213","track_id":"a5386297e19f1fa5cba3e6c745f88f7f0d09804a","url":"http://e.clickyab.com/show/sync/8bc27d90a0948342091dfa50cd5e253bcd62a6bc/35997/m184df33b0371e0cd392ee3307868fce792f52675c?parent=\u0026ref=\u0026s=3856691\u0026tid=53890855","width":300,"winner_cpm":5050},"fallback":"//a.clickyab.com/ads/show.php?a=2791492247287\u0026width=300\u0026height=250\u0026slot=92005149213\u0026domainname=entekhab.ir\u0026eventpage\u0026ck=true\u0026loc=http://entekhab.ir/\u0026ref=http://entekhab.ir/","height":250,"track_id":"92005149213","width":300}],"source":{"attributes":null,"floor_cpm":5050,"name":"entekhab.ir","soft_floor_cpm":6060,"supplier":{"excluded_demands":[],"floor_cpm":200,"name":"clickyab","share":1,"soft_floor_cpm":250}},"time":"2017-05-27T07:16:33.694113248Z","track_id":"8bc27d90a0948342091dfa50cd5e253bcd62a6bc","under_floor":false,"user_agent":"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"},"slot_id":"92005149213"}`
)

//func newImpression(t time.Time, slotCount int, source, sup string) exchange.Impression {
//	return mocks.Impression{
//		ITime: t,
//		ISource: mocks.Publisher{
//			PName: source,
//			PSupplier: mocks.Supplier{
//				SName: sup,
//			},
//		},
//		ISlots: make([]mocks.Slot, slotCount),
//	}
//}
//func newAdvertiser(cpm int64, dname string) exchange.Advertise {
//	return mocks.Advertiser{
//		MMaxCPM: cpm,
//		MDemand: mocks.Demand{
//			Mkey: dname,
//		},
//	}
//}

//func winnerToDelivery(imp exchange.Impression, ad exchange.Advertise, slot string) broker.Delivery {
//	job := materialize.WinnerJob(imp, ad, slot)
//	d, err := job.Encode()
//	assert.Nil(err)
//	return mocks.JsonDelivery{Data: d}
//}

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
	Convey("the demand test with the winner job", t, func() {
		//imp := newImpression(t1, 10, "test_winner", "test_demand")
		//adv := newAdvertiser(200, "adad")
		ctx, cl := context.WithCancel(base)
		defer cl()
		dem := consumer{ctx: ctx}
		delivery := dem.Consume()
		//data := winnerToDelivery(imp, adv, "aaa")
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
		So(t.AdOutCount, ShouldEqual, 1)
		So(t.AdOutBid, ShouldEqual, 5050)

		//So(t.Time, ShouldEqual, 1)
		//So(t.Request, ShouldEqual, 1)
		//So(t.Impression, ShouldEqual, 10)
		//
		//So(t.WinnerBid, ShouldBeZeroValue)
		//So(t.ShowBid, ShouldBeZeroValue)
		//So(t.Show, ShouldBeZeroValue)
		//So(t.Demand, ShouldBeZeroValue)
	})
}
