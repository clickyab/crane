package reducer

import (
	"context"
	"errors"
	"testing"
	"time"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/entity/mock_entity"
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
)

type filter struct {
}

var errTest = errors.New("test filter err")

func (filter) Check(_ entity.Context, c entity.Campaign) error {
	switch c.ID() {
	case 0:
		time.Sleep(time.Second)
		return nil
	case 1:
		return errTest
	case 2:
		return nil

	}
	panic("BUG")
}

func TestApply(t *testing.T) {
	convey.Convey("Filter should return ", t, func() {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		imp := mock_entity.NewMockContext(ctrl)
		pub := mock_entity.NewMockPublisher(ctrl)
		sup := mock_entity.NewMockSupplier(ctrl)

		sup.EXPECT().Name().Return("test").AnyTimes()
		pub.EXPECT().Supplier().Return(sup).AnyTimes()
		imp.EXPECT().Publisher().Return(pub).AnyTimes()
		convey.Convey("timeout error", func() {
			crt := mock_entity.NewMockCampaign(ctrl)
			crt.EXPECT().ID().Return(int32(0)).AnyTimes()

			res, err := Apply(ctx, imp, map[int32]entity.Campaign{1: crt}, []Filter{filter{}})
			convey.So(err, convey.ShouldBeError, ErrorTimeOut)
			convey.So(res, convey.ShouldBeNil)
		})
		convey.Convey("filter error", func() {
			crt := mock_entity.NewMockCampaign(ctrl)
			crt.EXPECT().ID().Return(int32(1)).AnyTimes()
			res, err := Apply(ctx, imp, map[int32]entity.Campaign{1: crt}, []Filter{filter{}})
			convey.So(err, convey.ShouldEqual, errTest)
			convey.So(res, convey.ShouldBeNil)
		})

		convey.Convey("creative", func() {
			crt := mock_entity.NewMockCampaign(ctrl)
			crt.EXPECT().ID().Return(int32(2)).AnyTimes()
			res, err := Apply(ctx, imp, map[int32]entity.Campaign{1: crt}, []Filter{filter{}})
			convey.So(err, convey.ShouldBeNil)
			convey.So(len(res), convey.ShouldEqual, 1)
		})
	})

}
