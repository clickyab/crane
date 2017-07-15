package native

import (
	"errors"
	"testing"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/entity/mock_entity"
	"clickyab.com/crane/crane/models/query"
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/goconvey/convey"
)

func TestImpression(t *testing.T) {
	c := gomock.NewController(t)
	convey.Convey("Native Should return ", t, func() {
		q := mock_entity.NewMockQPublisher(c)
		m := mock_entity.NewMockRequest(c)
		//		convey.Convey("error if domain is not set ", func() {
		//			m.EXPECT().Attributes().Return(map[string]string{})
		//			x, e := New(m)
		//			convey.So(x, convey.ShouldBeNil)
		//			convey.So(e, should.Equal, ErrorDomainIsEmplty)
		//		})
		convey.Convey("error if query return error ", func() {

			m.EXPECT().Attributes().Return(map[string]string{
				domain:   "d",
				supplier: "s",
			}).AnyTimes()
			q.EXPECT().ByPlatform(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, errors.New("db error"))
				//Do(func() (entity.Publisher, error) {
				//	return nil, errors.New("db error")
				//})
			query.Register(ll{})

			x, e := New(m)
			convey.So(x, convey.ShouldBeNil)
			convey.So(e, should.Equal, ErrorPublisherNotFound)
		})

	})
}

type ll struct {
}

func (ll) Find(int64) (entity.Publisher, error) {
	panic("implement me")
}

func (ll) ByPlatform(string, entity.Platforms, string) (entity.Publisher, error) {
	return nil, nil
}
