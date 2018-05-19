package filter_test

import (
	"testing"

	"clickyab.com/crane/demand/entity/mock_entity"
	"clickyab.com/crane/demand/filter"
	"github.com/golang/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDesktop_Check(t *testing.T) {
	Convey("checking app carrier filter", t, func() {
		ct := gomock.NewController(t)
		desktopSt := filter.Desktop{}
		context := mock_entity.NewMockContext(ct)
		creative := mock_entity.NewMockCreative(ct)
		campaign := mock_entity.NewMockCampaign(ct)
		creative.EXPECT().Campaign().Return(campaign).AnyTimes()

		Convey("request is mobile and campaign webmobile is off", func() {
			context.EXPECT().IsMobile().Return(true).AnyTimes()
			campaign.EXPECT().WebMobile().Return(false).AnyTimes()
			So(desktopSt.Check(context, creative), ShouldNotBeNil)
		})

		Convey("request is mobile and campaign webmobile is on", func() {
			context.EXPECT().IsMobile().Return(true).AnyTimes()
			campaign.EXPECT().WebMobile().Return(true).AnyTimes()
			So(desktopSt.Check(context, creative), ShouldBeNil)
		})

		Convey("request is not mobile and campaign web is on", func() {
			context.EXPECT().IsMobile().Return(false).AnyTimes()
			campaign.EXPECT().Web().Return(true).AnyTimes()
			So(desktopSt.Check(context, creative), ShouldBeNil)
		})

		Convey("request is not mobile and campaign web is off", func() {
			context.EXPECT().IsMobile().Return(false).AnyTimes()
			campaign.EXPECT().Web().Return(false).AnyTimes()
			So(desktopSt.Check(context, creative), ShouldNotBeNil)
		})

	})
}
