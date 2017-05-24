package rtb

import (
	"strconv"
	"testing"

	"clickyab.com/exchange/octopus/exchange"
	"clickyab.com/exchange/services/random"

	x "github.com/smartystreets/assertions"

	"clickyab.com/exchange/services/broker"
	"clickyab.com/exchange/services/broker/mock"

	"clickyab.com/exchange/octopus/exchange/mock_exchange"

	"clickyab.com/exchange/services/dset"
	dsetm "clickyab.com/exchange/services/dset/mock"

	"clickyab.com/exchange/services/dlock"
	dlockm "clickyab.com/exchange/services/dlock/mock"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

type situation struct {
	Case          int
	SCPMFloor     int64
	CPMFloor      int64
	MarginPercent int
	DSCPMFloor    int64
	DCPMFloor     int64
	UnderFloor    bool
	Demands       []dem
	WinnerDemand  int64
	Profit3rdAd   float64
	ProfitSuplier float64
	description   string
	Expectation   func(interface{}, ...interface{}) string
}

type dem struct {
	MaxCPM   int64
	HandyCap int64
}

func cases() []situation {
	return []situation{
		{1, 300, 200, 10, 330, 220, false, []dem{{470, 100}, {440, 110}}, 440, 44, 396, "", x.ShouldEqual},
		{2, 300, 200, 10, 330, 220, true, []dem{{470, 100}, {440, 110}}, 440, 44, 396, "", x.ShouldEqual},
		{3, 300, 200, 10, 330, 220, false, []dem{{470, 100}, {440, 100}}, 441, 44.1, 396.9, "", x.ShouldEqual},
		{4, 300, 200, 10, 330, 220, true, []dem{{470, 100}, {440, 100}}, 441, 44.1, 396.9, "", x.ShouldEqual},
		{5, 300, 200, 10, 330, 220, false, []dem{{470, 120}, {230, 115}}, 331, 33, 297, "", x.ShouldEqual},
		{6, 300, 200, 10, 330, 220, true, []dem{{470, 120}, {230, 115}}, 331, 33, 297, "", x.ShouldEqual},
		{7, 300, 200, 10, 330, 220, false, []dem{{340, 100}, {300, 125}}, 300, 30, 270, "", x.ShouldEqual},
		{8, 300, 200, 10, 330, 220, true, []dem{{340, 100}, {300, 125}}, 300, 30, 270, "", x.ShouldEqual},
		{9, 300, 200, 10, 330, 220, false, []dem{{230, 100}, {250, 110}}, 231, 23.1, 207.9, "", x.ShouldEqual},
		{10, 300, 200, 10, 330, 220, true, []dem{{230, 100}, {250, 110}}, 231, 23.1, 207.9, "", x.ShouldEqual},
		{11, 300, 200, 10, 330, 220, false, []dem{{200, 120}, {230, 90}}, 200, 22, 198, "", x.ShouldEqual},
		{12, 300, 200, 10, 330, 220, true, []dem{{200, 120}, {230, 90}}, 200, 20, 180, "", x.ShouldEqual},
		{14, 300, 200, 10, 330, 220, true, []dem{{190, 100}, {185, 105}}, 185, 18.5, 166.5, "", x.ShouldEqual},
		{15, 300, 200, 10, 330, 220, false, []dem{{250, 100}, {250, 100}}, 250, 25, 225, "", x.ShouldEqual},
		{16, 300, 200, 10, 330, 220, false, []dem{{350, 100}}, 331, 33, 297, "", x.ShouldEqual},
		{17, 300, 200, 10, 330, 220, true, []dem{{350, 100}}, 331, 33, 297, "", x.ShouldEqual},
		{18, 300, 200, 10, 330, 220, false, []dem{{310, 100}}, 221, 22, 198, "", x.ShouldEqual},
		{19, 300, 200, 10, 330, 220, true, []dem{{310, 100}}, 221, 22, 198, "", x.ShouldEqual},
		{20, 300, 200, 10, 330, 220, false, []dem{{240, 120}}, 221, 22, 198, "", x.ShouldEqual},
		{21, 300, 200, 10, 330, 220, true, []dem{{240, 120}}, 221, 22, 198, "", x.ShouldEqual},
		{23, 300, 200, 10, 330, 220, true, []dem{{210, 100}}, 210, 21, 189, "", x.ShouldEqual},
		{25, 300, 200, 10, 330, 220, true, []dem{{210, 80}}, 210, 21, 189, "", x.ShouldEqual},
		{27, 300, 200, 10, 330, 220, true, []dem{{210, 80}}, 210, 21, 189, "", x.ShouldEqual},
		{29, 300, 200, 10, 330, 220, true, []dem{{190, 100}}, 190, 21, 189, "", x.ShouldEqual},
	}

}

type Advertise struct {
	cpm,
	win int64
	demand exchange.Demand
}

func TestSelect(t *testing.T) {

	dset.Register(dsetm.NewMockDsetStore)
	dlock.Register(dlockm.NewMockDistributedLocker)

	ctrl := gomock.NewController(t)

	b := mock.GetChannelBroker()
	broker.SetActiveBroker(b)
	for _, u := range cases() {
		Convey("SelectCPM function test with case number "+strconv.Itoa(u.Case), t, func() {

			s := mock_exchange.NewMockSupplier(ctrl)
			s.EXPECT().Name().Return("test").AnyTimes()
			p := mock_exchange.NewMockPublisher(ctrl)
			p.EXPECT().FloorCPM().Return(u.DCPMFloor).AnyTimes()
			p.EXPECT().SoftFloorCPM().Return(u.DSCPMFloor).AnyTimes()
			p.EXPECT().Supplier().Return(s).AnyTimes()

			p.EXPECT().Rates().Return([]exchange.Rate{exchange.RateA}).AnyTimes()
			m := mock_exchange.NewMockImpression(ctrl)
			m.EXPECT().Source().Return(p).AnyTimes()

			m.EXPECT().PageTrackID().Return(<-random.ID).AnyTimes()
			m.EXPECT().UnderFloor().Return(u.UnderFloor)

			ads := make([]exchange.Advertise, 0)

			for _, a := range u.Demands {
				d := mock_exchange.NewMockDemand(ctrl)
				d.EXPECT().Handicap().Return(a.HandyCap).AnyTimes()
				ad := mock_exchange.NewMockAdvertise(ctrl)
				ad.EXPECT().MaxCPM().Return(a.MaxCPM).AnyTimes()
				ad.EXPECT().ID().Return(<-random.ID).AnyTimes()
				ad.EXPECT().Rates().Return([]exchange.Rate{exchange.RateB}).AnyTimes()
				ad.EXPECT().WinnerCPM().Return(u.WinnerDemand).AnyTimes()
				ad.EXPECT().SetWinnerCPM(gomock.Any()).AnyTimes()
				ad.EXPECT().Demand().Return(d).AnyTimes()
				ads = append(ads, ad)
			}

			slots := make(map[string][]exchange.Advertise)
			id := <-random.ID
			slots[id] = ads
			res := SelectCPM(m, slots)
			So(res[id].WinnerCPM(), u.Expectation, u.WinnerDemand)
		})
	}

	Convey("SelectCPM function Should return nil", t, func() {
		m := mock_exchange.NewMockImpression(ctrl)
		m.EXPECT().PageTrackID().Return(<-random.ID).AnyTimes()
		sup := mock_exchange.NewMockSupplier(ctrl)
		sup.EXPECT().Name().Return(<-random.ID).AnyTimes()
		p := mock_exchange.NewMockPublisher(ctrl)

		p.EXPECT().Supplier().Return(sup).AnyTimes()
		m.EXPECT().Source().Return(p).AnyTimes()
		id := <-random.ID
		s := make(map[string][]exchange.Advertise)
		s[id] = nil
		r := SelectCPM(m, s)
		So(r[id], ShouldBeNil)
	})

}
