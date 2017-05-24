package core

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"clickyab.com/exchange/octopus/exchange"
	mock_entity "clickyab.com/exchange/octopus/exchange/mock_exchange"
	"clickyab.com/exchange/services/random"

	"clickyab.com/exchange/services/config"

	"clickyab.com/exchange/services/broker"
	"clickyab.com/exchange/services/broker/mock"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/fzerorubigd/onion.v2"
)

func newPub(c *gomock.Controller) exchange.Publisher {
	s := mock_entity.NewMockSupplier(c)
	s.EXPECT().ExcludedDemands().Return([]string{}).AnyTimes()
	tmp := mock_entity.NewMockPublisher(c)
	tmp.EXPECT().Name().Return("publisher").AnyTimes()
	tmp.EXPECT().Supplier().Return(s).AnyTimes()
	tmp.EXPECT().FloorCPM().Return(int64(100)).AnyTimes()
	return tmp
}
func newImp(c *gomock.Controller, count int) exchange.Impression {
	tmp := make([]exchange.Slot, count)
	for i := range tmp {
		s := mock_entity.NewMockSlot(c)
		s.EXPECT().TrackID().Return(<-random.ID).AnyTimes()

		tmp[i] = s
	}
	l := mock_entity.NewMockLocation(c)
	l.EXPECT().Country().Return(exchange.Country{Name: "IRAN"}).AnyTimes()
	m := mock_entity.NewMockImpression(c)
	m.EXPECT().Scheme().Return("http").AnyTimes()
	m.EXPECT().Location().Return(l).AnyTimes()
	m.EXPECT().Slots().Return(tmp).AnyTimes()
	m.EXPECT().Source().Return(newPub(c)).AnyTimes()
	m.EXPECT().UnderFloor().Return(false).AnyTimes()

	return m
}

func TestProviders(t *testing.T) {
	def := onion.NewDefaultLayer()
	def.SetDefault("octupos.exchange.materialize.driver", "empty")
	config.Initialize("", "", "", def)
	ctrl := gomock.NewController(t)
	broker.SetActiveBroker(mock.GetChannelBroker())

	Convey("The provider's", t, func() {
		defer ctrl.Finish()
		maximumTimeout = 50 * time.Millisecond
		Reset(func() {
			allProviders = make(map[string]providerData)
		})

		Convey("Call func", func() {

			Convey("Should return two ads", func() {

				d1 := mock_entity.NewMockDemand(ctrl)
				d1.EXPECT().WhiteListCountries().Return([]string{}).AnyTimes()

				d1.EXPECT().Name().Return("d1").AnyTimes()

				d1.EXPECT().Handicap().Return(int64(100)).AnyTimes()
				d1.EXPECT().CallRate().Return(100).AnyTimes()
				d1.EXPECT().Provide(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
					Do(func(ctx context.Context, imp exchange.Impression, ch chan exchange.Advertise) {
						for _, s := range imp.Slots() {
							tmp := mock_entity.NewMockAdvertise(ctrl)
							tmp.EXPECT().MaxCPM().Return(int64(200)).AnyTimes()
							tmp.EXPECT().SlotTrackID().Return(s.TrackID()).AnyTimes()
							ch <- tmp
						}
						close(ch)
					}).AnyTimes()
				Register(d1, time.Millisecond*100)
				im := newImp(ctrl, 2)
				bk := context.Background()

				ads := Call(bk, im)
				So(len(ads), ShouldEqual, 2)
				So(len(ads[im.Slots()[0].TrackID()]), ShouldEqual, 1)
				So(len(ads[im.Slots()[1].TrackID()]), ShouldEqual, 1)

			})

			Convey("Should return NO ads", func() {
				d1 := mock_entity.NewMockDemand(ctrl)
				d1.EXPECT().WhiteListCountries().Return([]string{}).AnyTimes()
				d1.EXPECT().Name().Return("d1").AnyTimes()
				d1.EXPECT().Handicap().Return(int64(100)).AnyTimes()
				d1.EXPECT().CallRate().Return(100).AnyTimes()
				d1.EXPECT().Provide(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
					Do(func(ctx context.Context, imp exchange.Impression, ch chan exchange.Advertise) {
						time.Sleep(time.Millisecond * 150)
						for _, s := range imp.Slots() {
							tmp := mock_entity.NewMockAdvertise(ctrl)
							tmp.EXPECT().MaxCPM().Return(int64(200))
							tmp.EXPECT().SlotTrackID().Return(s.TrackID())
							ch <- tmp
						}
						close(ch)
					})
				Register(d1, time.Millisecond*100)
				im := newImp(ctrl, 2)
				bk := context.Background()

				ads := Call(bk, im)
				So(len(ads), ShouldEqual, 0)

			})

			Convey("Should return one provider with three ads (timeout test)", func() {
				d1 := mock_entity.NewMockDemand(ctrl)
				d1.EXPECT().WhiteListCountries().Return([]string{}).AnyTimes()
				d1.EXPECT().Name().Return("d1").AnyTimes()
				d1.EXPECT().Handicap().Return(int64(100)).AnyTimes()
				d1.EXPECT().CallRate().Return(100).AnyTimes()
				d1.EXPECT().Provide(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
					Do(func(ctx context.Context, imp exchange.Impression, ch chan exchange.Advertise) {
						time.Sleep(time.Millisecond * 100)
						for _, s := range imp.Slots() {
							tmp := mock_entity.NewMockAdvertise(ctrl)
							tmp.EXPECT().MaxCPM().Return(int64(200))
							tmp.EXPECT().SlotTrackID().Return(s.TrackID())
							ch <- tmp
						}
						close(ch)
					})
				Register(d1, time.Millisecond*100)
				d2 := mock_entity.NewMockDemand(ctrl)
				d2.EXPECT().WhiteListCountries().Return([]string{}).AnyTimes()
				d2.EXPECT().Name().Return("d2").AnyTimes()
				d2.EXPECT().Handicap().Return(int64(100)).AnyTimes()
				d2.EXPECT().CallRate().Return(100).AnyTimes()
				d2.EXPECT().Provide(gomock.Any(), gomock.Any(), gomock.Any()).
					Do(func(ctx context.Context, imp exchange.Impression, ch chan exchange.Advertise) {
						time.Sleep(time.Millisecond * 10)
						for _, s := range imp.Slots() {
							tmp := mock_entity.NewMockAdvertise(ctrl)
							tmp.EXPECT().MaxCPM().Return(int64(200))
							tmp.EXPECT().SlotTrackID().Return(s.TrackID())
							ch <- tmp
						}
						close(ch)
					})
				Register(d2, time.Millisecond*100)
				im := newImp(ctrl, 3)
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
				demand.EXPECT().Name().Return("test1").AnyTimes()
				Register(demand, time.Second*2)
				So(len(allProviders), ShouldEqual, 1)
				demand2 := mock_entity.NewMockDemand(ctrl)
				demand2.EXPECT().Name().Return("test1").AnyTimes()

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

	var counter [3000]int
	skips := [...]int{1, 10, 15, 27, 35, 48, 50, 68, 79, 87, 100}

	for _, s := range skips {
		Convey(fmt.Sprintf("Skip method should return true %d out of %d times hit for %d percent call rate.", int64(float64(len(counter))*(float64(s)/100.)), len(counter), s), t, func() {
			d := mock_entity.NewMockDemand(ctrl)
			d.EXPECT().CallRate().Return(s).AnyTimes()
			p := &providerData{name: <-random.ID, provider: d, timeout: time.Second}
			var tr int64
			wg := sync.WaitGroup{}
			for range counter {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if p.Skip() {
						atomic.AddInt64(&tr, 1)
					}
				}()
			}
			wg.Wait()
			So(tr, ShouldEqual, 3000-int64(float64(len(counter))*(float64(s)/100.)))

		})
	}

	Convey("Reset function should empty allProviders", t, func() {
		allProviders = make(map[string]providerData)
		allProviders["prv1"] = providerData{}
		allProviders["prv2"] = providerData{}
		ResetProviders()
		So(len(allProviders), ShouldEqual, 0)

	})

	Convey("Filters:", t, func() {

		Convey("isSameProvider function should return", func() {

			Convey("true if impression provider and provider are the same", func() {
				p2 := mock_entity.NewMockPublisher(ctrl)
				p2.EXPECT().Name().Return("prv1")
				m1 := mock_entity.NewMockImpression(ctrl)
				m1.EXPECT().Source().Return(p2)
				pd := providerData{name: "prv1"}
				So(isSameProvider(m1, pd), ShouldBeTrue)
			})

			Convey("false if impression provider and provider are NOT the same", func() {
				p1 := mock_entity.NewMockPublisher(ctrl)
				p1.EXPECT().Name().Return("prv1")
				m1 := mock_entity.NewMockImpression(ctrl)
				m1.EXPECT().Source().Return(p1)
				pd := providerData{name: "prv2"}
				So(isSameProvider(m1, pd), ShouldBeFalse)
			})
		})

		Convey("isNotwhitelistCountries function should return", func() {

			Convey("false if impression country is not in provider white list ", func() {

				pr := mock_entity.NewMockDemand(ctrl)
				pr.EXPECT().WhiteListCountries().Return([]string{"UAE", "IRAN"}).AnyTimes()
				pd := providerData{provider: pr}
				m := mock_entity.NewMockImpression(ctrl)
				l := mock_entity.NewMockLocation(ctrl)
				l.EXPECT().Country().Return(exchange.Country{ISO: "IRAN"})
				m.EXPECT().Location().Return(l)
				So(notWhitelistCountries(m, pd), ShouldBeFalse)
			})

			Convey("true if impression country is in provider white list ", func() {

				pr := mock_entity.NewMockDemand(ctrl)
				pr.EXPECT().WhiteListCountries().Return([]string{"UAE", "IRAN"}).AnyTimes()
				pd := providerData{provider: pr}
				m := mock_entity.NewMockImpression(ctrl)
				l := mock_entity.NewMockLocation(ctrl)
				l.EXPECT().Country().Return(exchange.Country{ISO: "USA"})
				m.EXPECT().Location().Return(l)
				So(notWhitelistCountries(m, pd), ShouldBeTrue)
			})
		})

		Convey("isExcludedDemands function should return", func() {

			Convey("true if impression exclude provider (by name)", func() {
				pub := mock_entity.NewMockPublisher(ctrl)
				sup := mock_entity.NewMockSupplier(ctrl)
				sup.EXPECT().ExcludedDemands().Return([]string{"PQ", "SAME", "PSD"})
				pub.EXPECT().Supplier().Return(sup)
				m := mock_entity.NewMockImpression(ctrl)
				m.EXPECT().Source().Return(pub)

				pd := providerData{name: "SAME"}
				So(isExcludedDemands(m, pd), ShouldBeTrue)

			})

			Convey("false if impression exclude provider (by name)", func() {
				pub := mock_entity.NewMockPublisher(ctrl)
				sup := mock_entity.NewMockSupplier(ctrl)
				sup.EXPECT().ExcludedDemands().Return([]string{"PQ", "EFG", "PSD"})
				pub.EXPECT().Supplier().Return(sup)
				m := mock_entity.NewMockImpression(ctrl)
				m.EXPECT().Source().Return(pub)

				pd := providerData{name: "UNIQUE"}

				So(isExcludedDemands(m, pd), ShouldBeFalse)
			})
		})
	})
}
