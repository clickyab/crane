package vast

import (
	"testing"

	"clickyab.com/crane/crane/entity/mock_entity"
	"clickyab.com/crane/crane/models/query"
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
)

func TestImpression(t *testing.T) {
	c := gomock.NewController(t)
	convey.Convey("test vast impression with no mode", t, func() {
		m := mock_entity.NewMockRequest(c)
		m.EXPECT().Attributes().Return(map[string]string{
			"l":        "long",
			"d":        "p30download.com",
			"supplier": "clickyab",
		}).AnyTimes()
		q := mock_entity.NewMockQPublisher(c)
		p := mock_entity.NewMockPublisher(c)
		q.EXPECT().ByPlatform(gomock.Any(), gomock.Any(), gomock.Any()).Return(p, nil).AnyTimes()
		p.EXPECT().UnderFloor().Return(true).AnyTimes()
		p.EXPECT().Name().Return("p30download.com").AnyTimes()
		p.EXPECT().Supplier().Return("clickyab").AnyTimes()
		p.EXPECT().FloorCPM().Return(int64(2000)).AnyTimes()
		p.EXPECT().SoftFloorCPM().Return(int64(3000)).AnyTimes()
		query.Register(q)
		imp, _ := New(m)
		convey.So(len(imp.Slots()), convey.ShouldEqual, 5)

	})
	convey.Convey("test with no break-type", t, func() {

		m := mock_entity.NewMockRequest(c)
		m.EXPECT().Attributes().Return(map[string]string{
			"l":        "invalid",
			"d":        "p30download.com",
			"supplier": "clickyab",
		}).AnyTimes()
		q := mock_entity.NewMockQPublisher(c)
		p := mock_entity.NewMockPublisher(c)
		q.EXPECT().ByPlatform(gomock.Any(), gomock.Any(), gomock.Any()).Return(p, nil).AnyTimes()
		p.EXPECT().UnderFloor().Return(true).AnyTimes()
		p.EXPECT().Name().Return("p30download.com").AnyTimes()
		p.EXPECT().Supplier().Return("clickyab").AnyTimes()
		p.EXPECT().FloorCPM().Return(int64(2000)).AnyTimes()
		p.EXPECT().SoftFloorCPM().Return(int64(3000)).AnyTimes()
		query.Register(q)
		imp, err := New(m)
		convey.So(err, convey.ShouldEqual, ErrorLenVast)
		convey.So(imp, convey.ShouldEqual, nil)
	})
	convey.Convey("test with no domain supplier", t, func() {

		m := mock_entity.NewMockRequest(c)
		m.EXPECT().Attributes().Return(map[string]string{
			"l": "short",
		}).AnyTimes()
		q := mock_entity.NewMockQPublisher(c)
		p := mock_entity.NewMockPublisher(c)
		q.EXPECT().ByPlatform(gomock.Any(), gomock.Any(), gomock.Any()).Return(p, nil).AnyTimes()
		p.EXPECT().UnderFloor().Return(true).AnyTimes()
		p.EXPECT().Name().Return("p30download.com").AnyTimes()
		p.EXPECT().Supplier().Return("clickyab").AnyTimes()
		p.EXPECT().FloorCPM().Return(int64(2000)).AnyTimes()
		p.EXPECT().SoftFloorCPM().Return(int64(3000)).AnyTimes()
		query.Register(q)
		imp, err := New(m)
		convey.So(err, convey.ShouldEqual, ErrorPublisherSupplierEmpty)
		convey.So(imp, convey.ShouldEqual, nil)
	})
}
