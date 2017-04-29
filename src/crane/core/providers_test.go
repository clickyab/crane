package core

import (
	"context"
	"crane/entity"
	"crane/entity/mock_entity"
	"net/http"
	"testing"
	"time"

	"net"

	. "github.com/smartystreets/goconvey/convey"

	"services/random"

	"github.com/golang/mock/gomock"
)

type slot struct {
	state       string
	size        int
	id          int64
	publicID    string
	ctr         float64
	allowedSize []int
	winnerBid   int64
	showURL     string
	advertise   entity.Advertise
}

func (hs *slot) ID() int64 {
	return hs.id
}

func (hs *slot) PublicID() string {
	return hs.publicID
}

func (hs *slot) SetShowURL(a string) {
	hs.showURL = a
}

func (hs *slot) SetWinnerAdvertise(a entity.Advertise) {
	hs.advertise = a
}

func (hs *slot) ShowURL() string {
	return hs.showURL
}

func (hs *slot) Size() int {
	return hs.size
}

func (hs *slot) SlotCTR() float64 {
	return hs.ctr
}

func (hs *slot) StateID() string {
	return hs.state
}

func (hs *slot) WinnerAdvertise() entity.Advertise {
	return hs.advertise
}

func (hs *slot) AllowedSize() []int {
	return hs.allowedSize
}

func (hs *slot) WinnerBid() int64 {
	return hs.winnerBid
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
		tmp[i] = slot{
			state: <-random.ID,
		}
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
				impression := newImp(2)
				bk := context.Background()

				ads := Call(bk, impression)
				So(len(ads), ShouldEqual, 2)
				So(len(ads[impression.Slots()[0].StateID()]), ShouldEqual, 1)
				So(len(ads[impression.Slots()[1].StateID()]), ShouldEqual, 1)

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
