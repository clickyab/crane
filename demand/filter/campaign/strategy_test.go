package campaign

import (
	"testing"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/entity/mock_entity"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStrategy_Check(t *testing.T) {
	Convey("test white list and black list publisher", t, func() {
		ct := gomock.NewController(t)
		strategySt := Strategy{}
		context := mock_entity.NewMockContext(ct)
		campaign := mock_entity.NewMockCampaign(ct)

		Convey("campaign is cpm but request is cpc", func() {
			campaign.EXPECT().Strategy().Return(entity.StrategyCPM).Times(2)
			context.EXPECT().Strategy().Return(entity.StrategyCPC).Times(2)
			So(strategySt.Check(context, campaign), ShouldNotBeNil)
		})

		Convey("both campaign and request is cpm", func() {
			campaign.EXPECT().Strategy().Return(entity.StrategyCPM).Times(1)
			context.EXPECT().Strategy().Return(entity.StrategyCPM).Times(1)
			So(strategySt.Check(context, campaign), ShouldBeNil)
		})
	})
}
