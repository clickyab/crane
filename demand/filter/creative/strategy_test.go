package creative

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
		creative := mock_entity.NewMockCreative(ct)
		campaign := mock_entity.NewMockCampaign(ct)
		creative.EXPECT().Campaign().Return(campaign).AnyTimes()

		Convey("campaign is cpm but request is cpc", func() {
			campaign.EXPECT().Strategy().Return(entity.StrategyCPM).Times(2)
			context.EXPECT().Strategy().Return(entity.StrategyCPC).Times(2)
			So(strategySt.Check(context, creative), ShouldNotBeNil)
		})

		Convey("both campaign and request is cpm", func() {
			campaign.EXPECT().Strategy().Return(entity.StrategyCPM).Times(1)
			context.EXPECT().Strategy().Return(entity.StrategyCPM).Times(1)
			So(strategySt.Check(context, creative), ShouldBeNil)
		})
	})
}
