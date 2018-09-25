package filter_test

import (
	"testing"

	"clickyab.com/crane/demand/entity/mock_entity"
	"clickyab.com/crane/demand/filter"
	"clickyab.com/crane/openrtb/v2.5"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConnectionType_Check(t *testing.T) {
	Convey("checking connection provider filter", t, func() {
		ct := gomock.NewController(t)
		connSt := filter.ConnectionType{}
		context := mock_entity.NewMockContext(ct)
		creative := mock_entity.NewMockCreative(ct)
		campaign := mock_entity.NewMockCampaign(ct)
		creative.EXPECT().Campaign().Return(campaign).AnyTimes()

		Convey("campaign has no conn type", func() {
			context.EXPECT().ConnectionType().Return(openrtb.ConnectionType_CELLULAT_NETWORK_2G).AnyTimes()
			campaign.EXPECT().ConnectionType().Return([]int{}).AnyTimes()
			So(connSt.Check(context, creative), ShouldBeNil)
		})

		Convey("campaign has conn type and did  match", func() {
			context.EXPECT().ConnectionType().Return(openrtb.ConnectionType_CELLULAT_NETWORK_2G).AnyTimes()
			campaign.EXPECT().ConnectionType().Return([]int{2, 4}).AnyTimes()
			So(connSt.Check(context, creative), ShouldBeNil)
		})

		Convey("campaign has conn type and did not match", func() {
			context.EXPECT().ConnectionType().Return(openrtb.ConnectionType_CELLULAT_NETWORK_3G).AnyTimes()
			campaign.EXPECT().ConnectionType().Return([]int{2, 4}).AnyTimes()
			So(connSt.Check(context, creative), ShouldNotBeNil)
		})

	})
}
