package core

import (
	"context"
	"entity"
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"net"

	"services/random"

	"entity/mock_entity"

	"github.com/golang/mock/gomock"
)

type slot struct {
	width, height int
	track         string
	max           int64
}

func (s slot) Width() int {
	return s.width
}

func (s slot) Height() int {
	return s.height
}

func (s slot) TrackID() string {
	return s.track
}

func (s slot) MaxCPM() int64 {
	return s.max
}

type imp struct {
	track  string
	ip     string
	ua     string
	undser bool

	slots []slot
	pub   entity.Publisher
	ts    *testing.T
}

func (imp) Attributes() map[string]interface{} {
	return nil
}

func (i imp) TrackID() string {
	return i.track
}

func (i imp) IP() net.IP {
	return net.ParseIP(i.ip)
}

func (i imp) UserAgent() string {
	return i.ua
}

func (i imp) Source() entity.Publisher {
	if i.pub == nil {
		ctrl := gomock.NewController(i.ts)
		tmp := mock_entity.NewMockPublisher(ctrl)
		tmp.EXPECT().FloorCPM().Return(int64(100))
		i.pub = tmp
	}
	return i.pub
}

func (i imp) Location() entity.Location {
	panic("implement me")
}

func (i imp) Slots() []entity.Slot {
	res := make([]entity.Slot, len(i.slots))
	for k := range i.slots {
		res[k] = &i.slots[k]
	}

	return res
}

func (i imp) Category() []entity.Category {
	panic("implement me")
}

func (i imp) Type() entity.ImpressionType {
	panic("implement me")
}

func (i imp) UnderFloor() bool {
	return i.undser
}

func (i imp) Raw() interface{} {
	panic("implement me")
}

func newImp(ts *testing.T, slotCount int) entity.Impression {
	tmp := make([]slot, slotCount)
	for i := range tmp {
		tmp[i] = slot{
			track: <-random.ID,
		}
	}
	return &imp{slots: tmp, ts: ts}
}

type tdemand struct {
	ts    *testing.T
	sleep time.Duration
	name  string
}

func (d *tdemand) Name() string {
	return d.name
}

func (*tdemand) Win(context.Context, string, int64) {
	panic("implement me")
}

func (*tdemand) Handicap() int64 {
	return 100
}

func (*tdemand) CallRate() int {
	return 100
}

func (d *tdemand) Status(ctx context.Context, rw http.ResponseWriter, rq *http.Request) {

}

func (d *tdemand) Provide(ctx context.Context, imp entity.Impression, ch chan map[string]entity.Advertise) {
	ctrl := gomock.NewController(d.ts)

	time.Sleep(d.sleep)
	ads := make(map[string]entity.Advertise)

	for _, s := range imp.Slots() {
		tmp := mock_entity.NewMockAdvertise(ctrl)
		tmp.EXPECT().MaxCPM().Return(int64(200))
		ads[s.TrackID()] = tmp
		ch <- ads
	}
	close(ch)
}

func TestProviders(t *testing.T) {
	ctrl := gomock.NewController(t)

	Convey("The provider's", t, func() {
		defer ctrl.Finish()
		maximumTimeout = 50 * time.Millisecond
		Reset(func() {
			allProviders = make(map[string]providerData)
		})

		Convey("Call func", func() {

			Convey("Should return two ads", func() {
				demand := &tdemand{t, time.Millisecond * 1, "test1"}
				Register(demand, time.Millisecond*100)
				im := newImp(t, 2)
				bk := context.Background()

				ads := Call(bk, im)
				So(len(ads), ShouldEqual, 2)
				So(len(ads[im.Slots()[0].TrackID()]), ShouldEqual, 1)
				So(len(ads[im.Slots()[1].TrackID()]), ShouldEqual, 1)

			})

			Convey("Should return NO ads", func() {
				demand := &tdemand{t, time.Millisecond * 100, "test1"}

				Register(demand, time.Millisecond*100)
				im := newImp(t, 2)
				bk := context.Background()

				ads := Call(bk, im)
				So(len(ads), ShouldEqual, 0)

			})

			Convey("Should return one provider with three ads (timeout test)", func() {
				demand1 := &tdemand{t, time.Millisecond * 100, "prv1"}
				Register(demand1, time.Millisecond*100)
				demand2 := &tdemand{t, time.Millisecond * 10, "prv2"}
				Register(demand2, time.Millisecond*100)
				im := newImp(t, 3)
				bk := context.Background()

				ads := Call(bk, im)
				So(len(ads), ShouldEqual, 3)
				So(len(ads[im.Slots()[0].TrackID()]), ShouldEqual, 1)
				So(len(ads[im.Slots()[1].TrackID()]), ShouldEqual, 1)
				So(len(ads[im.Slots()[2].TrackID()]), ShouldEqual, 1)

			})

		})

		Convey("Register func", func() {

			Convey("should panic if provider (name) is NOT unique", func() {
				demand := mock_entity.NewMockDemand(ctrl)
				demand.EXPECT().Name().Return("test1")
				Register(demand, time.Second*2)
				So(len(allProviders), ShouldEqual, 1)
				demand2 := mock_entity.NewMockDemand(ctrl)
				demand2.EXPECT().Name().Return("test1")

				So(func() {
					Register(demand2, time.Second*2)
				}, ShouldPanic)

			})

			Convey("should register multiple providers", func() {
				demand := mock_entity.NewMockDemand(ctrl)
				demand.EXPECT().Name().Return("test1")

				Register(demand, time.Second*2)
				So(len(allProviders), ShouldEqual, 1)
				demand2 := mock_entity.NewMockDemand(ctrl)
				demand2.EXPECT().Name().Return("test2")
				So(func() {
					Register(
						demand2, time.Second*2)
				}, ShouldNotPanic)
				So(len(allProviders), ShouldEqual, 2)

			})

		})
	})
}
