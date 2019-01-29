package campaign

import (
	"testing"

	"clickyab.com/crane/demand/entity/mock_entity"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAppCarrier_Check(t *testing.T) {
	Convey("checking app carrier filter", t, func() {
		ct := gomock.NewController(t)
		carrierSt := AppCarrier{}
		context := mock_entity.NewMockContext(ct)
		creative := mock_entity.NewMockCreative(ct)
		campaign := mock_entity.NewMockCampaign(ct)
		creative.EXPECT().Campaign().Return(campaign).AnyTimes()

		Convey("campaign have no carriers", func() {
			contextCarriers := "some carrier"
			context.EXPECT().Carrier().Return(contextCarriers).Times(1)
			campaign.EXPECT().AppCarriers().Return([]string{}).Times(1)
			So(carrierSt.Check(context, campaign), ShouldBeNil)

		})

		Convey("campaign carriers didn't match the context carrier", func() {
			contextCarriers := "some carrier"
			context.EXPECT().Carrier().Return(contextCarriers).Times(1)
			campaign.EXPECT().AppCarriers().Return([]string{"some other carrier"}).Times(1)
			So(carrierSt.Check(context, campaign), ShouldNotBeNil)
		})

	})
}
