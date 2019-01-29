package creative

import (
	"testing"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/entity/mock_entity"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestISP_Check(t *testing.T) {
	Convey("checking app carrier filter", t, func() {
		ct := gomock.NewController(t)
		ispSt := ISP{}
		context := mock_entity.NewMockContext(ct)
		creative := mock_entity.NewMockCreative(ct)
		campaign := mock_entity.NewMockCampaign(ct)
		location := mock_entity.NewMockLocation(ct)
		context.EXPECT().Location().Return(location).AnyTimes()
		creative.EXPECT().Campaign().Return(campaign).AnyTimes()

		Convey("campaign have no isp", func() {
			campaign.EXPECT().ISP().Return([]string{}).AnyTimes()
			var contextIsp = entity.ISP{ID: 4}
			location.EXPECT().ISP().Return(contextIsp).Times(1)
			So(ispSt.Check(context, creative), ShouldBeNil)

		})

		Convey("campaign have isps but not match the context", func() {
			campIsp := []string{"1", "2", "3"}
			campaign.EXPECT().ISP().Return(campIsp).AnyTimes()
			var contextIsp = entity.ISP{ID: 4}
			location.EXPECT().ISP().Return(contextIsp).Times(1)
			So(ispSt.Check(context, creative), ShouldNotBeNil)
		})

		Convey("campaign have isps and match the context", func() {
			campIsp := []string{"1", "2", "3"}
			campaign.EXPECT().ISP().Return(campIsp).AnyTimes()
			var contextIsp = entity.ISP{ID: 1}
			location.EXPECT().ISP().Return(contextIsp).Times(1)
			So(ispSt.Check(context, creative), ShouldBeNil)
		})

	})
}
