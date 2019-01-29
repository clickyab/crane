package campaign

import (
	"testing"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/entity/mock_entity"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOS_Check(t *testing.T) {
	Convey("checking os filter", t, func() {
		ct := gomock.NewController(t)
		osSt := OS{}
		context := mock_entity.NewMockContext(ct)
		campaign := mock_entity.NewMockCampaign(ct)

		Convey("campaign has no os", func() {
			campaign.EXPECT().AllowedOS().Return([]string{}).AnyTimes()
			So(osSt.Check(context, campaign), ShouldBeNil)
		})

		Convey("campaign os not empty but context os not valid", func() {
			contextOs := entity.OS{Valid: false}
			campaign.EXPECT().AllowedOS().Return([]string{"1", "2"}).AnyTimes()
			context.EXPECT().OS().Return(contextOs).AnyTimes()
			So(osSt.Check(context, campaign), ShouldNotBeNil)
		})

		Convey("campaign os not empty and context not empty and it match", func() {
			contextOs := entity.OS{Valid: true, ID: 1}
			campaign.EXPECT().AllowedOS().Return([]string{"1", "2"}).AnyTimes()
			context.EXPECT().OS().Return(contextOs).AnyTimes()
			So(osSt.Check(context, campaign), ShouldBeNil)
		})

		Convey("campaign os not empty and context not empty and it not match", func() {
			contextOs := entity.OS{Valid: true, ID: 3}
			campaign.EXPECT().AllowedOS().Return([]string{"1", "2"}).AnyTimes()
			context.EXPECT().OS().Return(contextOs).AnyTimes()
			So(osSt.Check(context, campaign), ShouldNotBeNil)
		})

	})
}
