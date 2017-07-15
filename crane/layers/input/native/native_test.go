package native

import (
	"errors"
	"testing"

	"clickyab.com/crane/crane/entity/mock_entity"
	"clickyab.com/crane/crane/models/query"
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/goconvey/convey"
)

func TestImpression(t *testing.T) {
	c := gomock.NewController(t)
	convey.Convey("Native Should return ", t, func() {

		m := mock_entity.NewMockRequest(c)

		convey.Convey("error if domain is not set ", func() {
			m.EXPECT().Attributes().Return(map[string]string{})

			x, e := New(m)

			convey.So(x, convey.ShouldBeNil)
			convey.So(e, should.Equal, ErrorDomainIsEmplty)

		})

		convey.Convey("error if count is not valid ", func() {

			m.EXPECT().Attributes().Return(map[string]string{
				domain:   "d",
				supplier: "s",
			}).AnyTimes()

			x, e := New(m)

			convey.So(x, convey.ShouldBeNil)
			convey.So(e, should.Equal, ErrorCountNotValid)

		})

		convey.Convey("error if query ( DB ) return error ", func() {

			m.EXPECT().Attributes().Return(map[string]string{
				domain:   "d",
				supplier: "s",
				count:    "4",
			}).AnyTimes()

			q := mock_entity.NewMockQPublisher(c)
			q.EXPECT().ByPlatform(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, errors.New("db error"))
			query.Register(q)

			x, e := New(m)

			convey.So(x, convey.ShouldBeNil)
			convey.So(e, should.Equal, ErrorPublisherNotFound)

		})

		convey.Convey("impression ", func() {

			m.EXPECT().Attributes().Return(map[string]string{
				domain:   "d",
				supplier: "s",
				count:    "4",
			}).AnyTimes()

			p := mock_entity.NewMockPublisher(c)
			p.EXPECT().UnderFloor().Return(true)
			p.EXPECT().Name().Return("publisher").AnyTimes()
			p.EXPECT().Supplier().Return("supplier").AnyTimes()
			p.EXPECT().FloorCPM().Return(int64(1))
			p.EXPECT().SoftFloorCPM().Return(int64(1))
			p.EXPECT().UnderFloor().Return(true)

			q := mock_entity.NewMockQPublisher(c)
			q.EXPECT().ByPlatform(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(p, nil)
			query.Register(q)

			x, e := New(m)

			convey.So(x, convey.ShouldNotBeNil)
			convey.So(len(x.Slots()), convey.ShouldEqual, 4)
			convey.So(e, convey.ShouldBeNil)

		})
	})
}
