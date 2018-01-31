package capping

import (
	"testing"

	"sort"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/entity/mock_entity"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/kv/mock"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func creatives(ct *gomock.Controller) []entity.Creative {
	res := make([]entity.Creative, 0)
	for i := 0; i < 100; i++ {
		cr := mock_entity.NewMockCreative(ct)
		cp := mock_entity.NewMockCampaign(ct)
		cp.EXPECT().Frequency().Return(2).AnyTimes()
		cp.EXPECT().ID().Return(int64(i)).AnyTimes()

		cr.EXPECT().Campaign().Return(cp).AnyTimes()
		cr.EXPECT().Size().Return(1).AnyTimes()
		cr.EXPECT().ID().Return(int64(i)).AnyTimes()
		cr.EXPECT().SetCapping(gomock.Any()).Do(func(b entity.Capping) {
			cr.EXPECT().Capping().Return(b).AnyTimes()
		}).AnyTimes()

		res = append(res, cr)
	}
	return res
}

func TestApplyCapping(t *testing.T) {
	Convey("Capping apply ", t, func() {
		ct := gomock.NewController(t)
		cx := mock_entity.NewMockContext(ct)
		user := mock_entity.NewMockUser(ct)
		user.EXPECT().ID().Return("one").AnyTimes()
		cx.EXPECT().User().Return(user).AnyTimes()
		cx.EXPECT().EventPage().Return("eventPage").AnyTimes()

		Convey("none mode", func() {
			var crs = creatives(ct)

			res := ApplyCapping(entity.CappingNone, "none", crs, "eventPage")
			So(len(res), ShouldEqual, len(crs))
			for i := range res {
				res[i].Capping().Store(res[i].ID())
			}

			for i := 0; i < 24; i++ {
				res[i].Capping().Store(res[i].ID())
			}

			for i := 0; i < 56; i++ {
				res[i].Capping().Store(res[i].ID())
			}

			crs = creatives(ct)

			ress := ApplyCapping(entity.CappingNone, "none", crs, "eventPage")

			So(len(ress), ShouldEqual, 100)
		})

		Convey("strict mode", func() {
			var crs = creatives(ct)

			res := ApplyCapping(entity.CappingStrict, "strict", crs, "eventPage")
			So(len(res), ShouldEqual, len(crs))
			for i := range res {
				res[i].Capping().Store(res[i].ID())
			}

			for i := 0; i < 24; i++ {
				res[i].Capping().Store(res[i].ID())
			}

			for i := 0; i < 56; i++ {
				res[i].Capping().Store(res[i].ID())
			}

			crs = creatives(ct)

			xres := ApplyCapping(entity.CappingStrict, "strict", crs, "eventPage")
			So(len(xres), ShouldEqual, 76)
			s := sortByCap(xres)

			sort.Sort(s)
			a, b := s[:44], s[44:]
			for _, x := range a {
				So(x.Capping().View(), ShouldEqual, 1)
			}

			for _, x := range b {
				So(x.Capping().View(), ShouldEqual, 2)
			}
		})

		Convey("reset mode", func() {
			var crs = creatives(ct)

			res := ApplyCapping(entity.CappingReset, "reset", crs, "eventPage")
			So(len(res), ShouldEqual, len(crs))
			for i := range res {
				res[i].Capping().Store(res[i].ID())
			}

			for i := 0; i < 24; i++ {
				res[i].Capping().Store(res[i].ID())
			}

			for i := 0; i < 56; i++ {
				res[i].Capping().Store(res[i].ID())
			}

			crs = creatives(ct)
			xres := ApplyCapping(entity.CappingReset, "reset", crs, "eventPage")
			s := sortByCap(xres)
			sort.Sort(s)
			So(len(s), ShouldEqual, 100)
			a, b, c := s[:44], s[44:76], s[76:]
			for _, x := range a {
				So(x.Capping().View(), ShouldEqual, 1)
			}

			for _, x := range b {
				So(x.Capping().View(), ShouldEqual, 2)
			}

			for _, x := range c {
				So(x.Capping().View(), ShouldEqual, 3)
			}

			var views = 0
			for i := range s {
				views += s[i].Capping().View()
			}
			So(views, ShouldEqual, 180)
		})
	})
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
