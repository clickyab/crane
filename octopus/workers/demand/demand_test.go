package demand

//import . "github.com/smartystreets/goconvey/convey"

//var (
//	t1, _ = time.Parse("2006-01-02T15:04:05.000Z", "2017-03-21T00:01:00.000Z")
//	t2, _ = time.Parse("2006-01-02T15:04:05.000Z", "2017-03-21T01:01:00.000Z")
//	t3, _ = time.Parse("2006-01-02T15:04:05.000Z", "2017-03-21T02:01:00.000Z")
//)
//
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
//
//func demandToDelivery(in exchange.Impression, dmn exchange.Demand, ads map[string]exchange.Advertise) broker.Delivery {
//	job := materialize.DemandJob(in, dmn, ads)
//	d, err := job.Encode()
//	assert.Nil(err)
//	return mocks.JsonDelivery{Data: d}
//}
//
//type agg struct {
//	c chan datamodels.TableModel
//}
//
//func (a *agg) Channel() chan<- datamodels.TableModel {
//	return a.c
//}
//
//func TestDemand(t *testing.T) {
//	config.Initialize("test", "test", "test")
//	// make sure channel has space for more than 1 delivery
//	a := &agg{c: make(chan datamodels.TableModel, 2)}
//	datamodels.RegisterAggregator(a)
//	//base := context.Background()
//	Convey("the demand test with the impression job", t, func() {
//
//	})
//}
