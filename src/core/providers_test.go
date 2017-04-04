package core

import (
	"context"
	"entity"
	"entity/mock_entity"
	"net/http"
	"testing"
	"time"

	"net"

	"services/random"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/golang/mock/gomock"
)

type slot struct {
	state string
}

func (*slot) ID() int64 {
	panic("implement me")
}

func (*slot) PublicID() string {
	panic("implement me")
}

func (*slot) Size() int {
	panic("implement me")
}

func (s *slot) StateID() string {
	if s.state == "" {
		s.state = <-random.ID
	}

	return s.state
}

func (*slot) SlotCTR() float64 {
	panic("implement me")
}

func (*slot) SetWinnerAdvertise(entity.Advertise) {
	panic("implement me")
}

func (*slot) WinnerAdvertise() entity.Advertise {
	panic("implement me")
}

func (*slot) SetShowURL(string) {
	panic("implement me")
}

type imp struct {
	slots []slot
}

func (*imp) Attributes(entity.ImpressionAttributes) interface{} {
	panic("implement me")
}

func (*imp) Category() []entity.Category {
	panic("implement me")
}

func (*imp) ClientID() int64 {
	panic("implement me")
}

func (*imp) IP() net.IP {
	panic("implement me")
}

func (*imp) Location() entity.Location {
	panic("implement me")
}

func (*imp) MegaIMP() string {
	panic("implement me")
}

func (*imp) OS() entity.OS {
	panic("implement me")
}

func (*imp) Request() *http.Request {
	panic("implement me")
}

func (i *imp) Slots() []entity.Slot {
	tmp := make([]entity.Slot, len(i.slots))
	for j := range i.slots {
		tmp[j] = &i.slots[j]
	}
	return tmp
}

func (*imp) Source() entity.Publisher {
	panic("implement me")
}

func (*imp) UserAgent() string {
	panic("implement me")
}

func newImp(slotCount int) entity.Impression {
	tmp := make([]slot, slotCount)
	for i := range tmp {
		tmp[i] = slot{}
	}
	return &imp{tmp}
}

type tdemand struct {
	ts    *testing.T
	sleep time.Duration
}

func (d *tdemand) Status(ctx context.Context, rw http.ResponseWriter, rq *http.Request) {

}

func (d *tdemand) Provide(ctx context.Context, imp entity.Impression, ch chan map[string]entity.Advertise) {
	ctrl := gomock.NewController(d.ts)

	time.Sleep(d.sleep)
	ads := make(map[string]entity.Advertise)

	for _, s := range imp.Slots() {
		ads[s.StateID()] = mock_entity.NewMockAdvertise(ctrl)
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
			allProviders = allProviders[:0]
		})

		Convey("Call func", func() {

			Convey("Should return two ads", func() {
				demand := &tdemand{t, time.Millisecond * 1}
				Register("prv1", demand, time.Millisecond*100)
				im := newImp(2)
				bk := context.Background()

				ads := Call(bk, im)
				So(len(ads), ShouldEqual, 2)
				So(len(ads[im.Slots()[0].StateID()]), ShouldEqual, 1)
				So(len(ads[im.Slots()[1].StateID()]), ShouldEqual, 1)

			})

			Convey("Should return NO ads", func() {
				demand := &tdemand{t, time.Millisecond * 100}

				Register("prv1", demand, time.Millisecond*100)
				im := newImp(2)
				bk := context.Background()

				ads := Call(bk, im)
				So(len(ads), ShouldEqual, 0)

			})

			Convey("Should return one provider with three ads (timeout test)", func() {
				demand1 := &tdemand{t, time.Millisecond * 100}
				Register("prv1", demand1, time.Millisecond*100)
				demand2 := &tdemand{t, time.Millisecond * 10}
				Register("prv2", demand2, time.Millisecond*100)
				im := newImp(3)
				bk := context.Background()

				ads := Call(bk, im)
				So(len(ads), ShouldEqual, 3)
				So(len(ads[im.Slots()[0].StateID()]), ShouldEqual, 1)
				So(len(ads[im.Slots()[1].StateID()]), ShouldEqual, 1)
				So(len(ads[im.Slots()[2].StateID()]), ShouldEqual, 1)

			})

		})

		Convey("Register func", func() {

			Convey("should panic if provider (name) is NOT unique", func() {
				demand := mock_entity.NewMockDemand(ctrl)
				Register("First Provider", demand, time.Second*2)
				So(len(allProviders), ShouldEqual, 1)
				So(func() {
					Register("First Provider",
						demand, time.Second*2)
				}, ShouldPanic)

			})

			Convey("should register multiple providers", func() {
				demand := mock_entity.NewMockDemand(ctrl)
				Register("First Provider", demand, time.Second*2)
				So(len(allProviders), ShouldEqual, 1)
				So(func() {
					Register("Second Provider",
						demand, time.Second*2)
				}, ShouldNotPanic)
				So(len(allProviders), ShouldEqual, 2)

			})

		})
	})
}
