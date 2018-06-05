package rtb

import (
	"testing"

	"fmt"

	"math/rand"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/entity/mock_entity"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/kv/mock"
	"github.com/clickyab/services/random"
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
)

func TestSelect(t *testing.T) {
	for _, s := range samples(t) {
		convey.Convey(s.title, t, func() {
			internalSelect(s.context, s.ads)
			var winner entity.Creative
			var seat entity.Seat
			for _, ta := range s.context.Seats() {
				if x := ta.WinnerAdvertise(); x != nil {
					if winner != nil {
						t.Logf("[BUG] we have more then one winner!!")
						t.Fail()
						break
					} else {
						winner = x
						seat = ta
					}
				}
			}
			if winner == nil {
				convey.So(s.winner, convey.ShouldBeNil)
			} else {
				convey.So(seat.CPM(), convey.ShouldEqual, s.cpm)
				convey.So(seat.Bid(), convey.ShouldAlmostEqual, s.bid, 5)
				convey.So(winner.ID(), convey.ShouldEqual, s.winner.ID())

			}
		})
	}
}
func samples(t *testing.T) []sample {
	c := gomock.NewController(t)
	res := make([]sample, 0)

	{
		col1 := make([]entity.Creative, 0)
		ad1 := creativeMaker(400, .05, entity.StrategyCPM, c)
		ctx := contextMaker(500, nil, false, entity.StrategyCPM, c)

		col1 = append(col1, ad1)
		res = append(res, sample{
			title:   "rtb case one",
			winner:  nil,
			cpm:     0,
			bid:     0,
			context: ctx,
			ads:     col1,
		})
	}

	// {
	// 	col1 := make([]entity.Creative, 0)
	// 	ad1 := creativeMaker(16000, .05, entity.StrategyCPM, c)
	// 	ctx := contextMaker(500, ad1, false, entity.StrategyCPM, c)

	// 	col1 = append(col1, ad1)
	// 	res = append(res, sample{
	// 		title:   "rtb case two",
	// 		winner:  ad1,
	// 		cpm:     500,
	// 		bid:     500 / 0.75,
	// 		context: ctx,
	// 		ads:     col1,
	// 	})
	// }

	{
		col1 := make([]entity.Creative, 0)
		ad1, ad2 := creativeMaker(16000, .05, entity.StrategyCPM, c),
			creativeMaker(10000, .05, entity.StrategyCPM, c)
		ctx := contextMaker(500, ad1, false, entity.StrategyCPM, c)

		col1 = append(col1, ad1, ad2)
		res = append(res, sample{
			title:   "rtb case three",
			winner:  ad1,
			cpm:     10010,
			bid:     10010 / 0.75,
			context: ctx,
			ads:     col1,
		})
	}

	{
		col1 := make([]entity.Creative, 0)
		ad1, ad2, ad3 := creativeMaker(16000, .05, entity.StrategyCPM, c),
			creativeMaker(12000, .05, entity.StrategyCPM, c),
			creativeMaker(10000, .05, entity.StrategyCPM, c)

		ctx := contextMaker(500, ad1, false, entity.StrategyCPM, c)

		col1 = append(col1, ad1, ad2, ad3)
		res = append(res, sample{
			title:   "rtb case four",
			winner:  ad1,
			cpm:     12010,
			bid:     12010 / 0.75,
			context: ctx,
			ads:     col1,
		})
	}

	{
		col1 := make([]entity.Creative, 0)
		ad1 := creativeMaker(400, .05, entity.StrategyCPC, c)
		ctx := contextMaker(500, nil, false, entity.StrategyCPC, c)

		col1 = append(col1, ad1)
		res = append(res, sample{
			title:   "rtb case five",
			winner:  nil,
			cpm:     0,
			bid:     0,
			context: ctx,
			ads:     col1,
		})
	}

	// {
	// 	col1 := make([]entity.Creative, 0)
	// 	ad1 := creativeMaker(16000, .05, entity.StrategyCPC, c)
	// 	ctx := contextMaker(500, ad1, false, entity.StrategyCPC, c)

	// 	col1 = append(col1, ad1)
	// 	res = append(res, sample{
	// 		title:   "rtb case six",
	// 		winner:  ad1,
	// 		cpm:     500,
	// 		bid:     500 / 0.75,
	// 		context: ctx,
	// 		ads:     col1,
	// 	})
	// }

	{
		col1 := make([]entity.Creative, 0)
		ad1, ad2 := creativeMaker(16000, .05, entity.StrategyCPC, c),
			creativeMaker(10000, .05, entity.StrategyCPC, c)
		ctx := contextMaker(500, ad1, false, entity.StrategyCPC, c)

		col1 = append(col1, ad1, ad2)
		res = append(res, sample{
			title:   "rtb case seven",
			winner:  ad1,
			cpm:     7510,  //int64(math.Floor(10010.0 * 0.75)),
			bid:     10010, // its calculated stuff a little bit off ignore
			context: ctx,
			ads:     col1,
		})
	}

	{
		col1 := make([]entity.Creative, 0)
		ad1, ad2, ad3 := creativeMaker(16000, .05, entity.StrategyCPC, c),
			creativeMaker(12000, .05, entity.StrategyCPC, c),
			creativeMaker(10000, .05, entity.StrategyCPC, c)

		ctx := contextMaker(500, ad1, false, entity.StrategyCPC, c)

		col1 = append(col1, ad1, ad2, ad3)
		res = append(res, sample{
			title:   "rtb case eight",
			winner:  ad1,
			cpm:     9010,  // 12010 * 0.75,
			bid:     12010, // its calculated and its ok
			context: ctx,
			ads:     col1,
		})
	}

	return res
}

type sample struct {
	title   string
	ads     []entity.Creative
	bid     float64
	cpm     int64
	context entity.Context
	winner  entity.Creative
}

var default_floor_cpm = map[string]int64{"web_vast": 2000,
	"web_banner": 2000,
	"web_native": 2000,
	"app_vast":   2000,
	"app_native": 2000,
	"app_banner": 2000}

func contextMaker(minbid int64, winner entity.Creative, underfloor bool, strategy entity.Strategy, c *gomock.Controller) entity.Context {

	ctx := mock_entity.NewMockContext(c)
	ctx.EXPECT().EventPage().Return(<-random.ID).AnyTimes()
	user := mock_entity.NewMockUser(c)
	user.EXPECT().ID().Return(<-random.ID).AnyTimes()
	ctx.EXPECT().User().Return(user).AnyTimes()
	ctx.EXPECT().MultiVideo().Return(false).AnyTimes()
	ctx.EXPECT().Capping().Return(entity.CappingNone).AnyTimes()
	ctx.EXPECT().UnderFloor().Return(underfloor).AnyTimes()

	//TODO: make sample data for statistics
	ctx.EXPECT().GetCreativesStatistics().Return(make([]entity.CreativeStatistics, 0)).AnyTimes()
	seat := mock_entity.NewMockSeat(c)
	seat.EXPECT().RequestType().Return(entity.RequestTypeBanner).AnyTimes()
	seat.EXPECT().Type().Return(entity.InputTypeDemand).AnyTimes()
	seat.EXPECT().MinBid().Return(minbid).AnyTimes()
	seat.EXPECT().MinCPC().Return(float64(minbid)).AnyTimes()
	// TODO : fix this
	seat.EXPECT().MinCPM().Return(float64(0)).AnyTimes()
	seat.EXPECT().SoftCPM().Return(float64(0)).AnyTimes()
	seat.EXPECT().Acceptable(gomock.Any()).Return(true).AnyTimes()
	seat.EXPECT().PublicID().Return("213456").AnyTimes()
	seat.EXPECT().CTR().Return(.1).AnyTimes()
	seat.EXPECT().Size().Return(20).AnyTimes()
	if winner == nil {
		seat.EXPECT().WinnerAdvertise().Return(nil).AnyTimes()
	}

	seat.EXPECT().SetWinnerAdvertise(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(creative entity.Creative, bid, share float64) {
			seat.EXPECT().WinnerAdvertise().Return(creative).AnyTimes()

			seat.EXPECT().Bid().Return(bid).AnyTimes()
			seat.EXPECT().CPM().Return(share).AnyTimes()
		}).AnyTimes()
	ctx.EXPECT().Seats().Return([]entity.Seat{seat}).AnyTimes()

	sup := mock_entity.NewMockSupplier(c)
	sup.EXPECT().Share().Return(100).AnyTimes()
	sup.EXPECT().Strategy().Return(strategy).AnyTimes()
	sup.EXPECT().Name().Return("clickyab").AnyTimes()

	softCPC := sup.EXPECT().SoftFloorCPC(gomock.Any(), gomock.Any())
	softCPC.Do(func(a, b string) {
		softCPC.Return(default_floor_cpm[fmt.Sprintf("%s_%s", a, b)]).AnyTimes()
	}).AnyTimes()

	softCPM := sup.EXPECT().SoftFloorCPM(gomock.Any(), gomock.Any())
	softCPM.Do(func(a, b string) {
		softCPM.Return(default_floor_cpm[fmt.Sprintf("%s_%s", a, b)]).AnyTimes()
	}).AnyTimes()

	pub := mock_entity.NewMockPublisher(c)
	pub.EXPECT().ID().Return(int64(32456)).AnyTimes()
	pub.EXPECT().Supplier().Return(sup).AnyTimes()
	pub.EXPECT().Type().Return(entity.PublisherTypeWeb).AnyTimes()

	ctx.EXPECT().Publisher().Return(pub).AnyTimes()
	return ctx
}

func creativeMaker(maxbid int64, ctr float64, strategy entity.Strategy, c *gomock.Controller) entity.Creative {

	rnd := rand.Int63()
	ad := mock_entity.NewMockCreative(c)
	ad.EXPECT().ID().Return(rnd).AnyTimes()
	ad.EXPECT().MaxBID().Return(maxbid).AnyTimes()
	ad.EXPECT().Type().Return(entity.AdTypeBanner).AnyTimes()
	ad.EXPECT().AdCTR().Return(ctr).AnyTimes()
	ad.EXPECT().SetCapping(gomock.Any()).Do(func(capping entity.Capping) {
		ad.EXPECT().Capping().Return(capping).AnyTimes()
	})

	campaign := mock_entity.NewMockCampaign(c)
	campaign.EXPECT().ID().Return(rnd + 1).AnyTimes()
	campaign.EXPECT().Strategy().Return(strategy).AnyTimes()
	campaign.EXPECT().Frequency().Return(10).AnyTimes()
	ad.EXPECT().Campaign().Return(campaign).AnyTimes()
	return ad
}

func init() {
	kv.Register(mock.NewMockStore,
		mock.NewMockChannelStore,
		mock.NewMockDistributedLocker,
		mock.NewMockDsetStore,
		mock.NewAtomicMockStore,
		mock.NewCacheMock(),
		nil,
		kv.NewOneTimeSetter,
	)
}
