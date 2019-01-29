package campaign

import (
	"testing"

	"clickyab.com/crane/demand/entity/mock_entity"
	openrtb "clickyab.com/crane/openrtb/v2.5"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCategory_Check(t *testing.T) {
	Convey("test iab category filter", t, func() {
		ct := gomock.NewController(t)
		catSt := Category{}
		context := mock_entity.NewMockContext(ct)
		campaign := mock_entity.NewMockCampaign(ct)

		Convey("campaign category is empty", func() {
			campaign.EXPECT().Category().Return(make([]openrtb.ContentCategory, 0)).AnyTimes()
			So(catSt.Check(context, campaign), ShouldBeNil)
		})

		Convey("campaign category is not empty but request cat is empty", func() {
			campCat := []openrtb.ContentCategory{
				openrtb.ContentCategory_IAB1,
				openrtb.ContentCategory_IAB1S5,
			}
			campaign.EXPECT().Category().Return(campCat).AnyTimes()
			context.EXPECT().Category().Return(make([]openrtb.ContentCategory, 0)).AnyTimes()
			So(catSt.Check(context, campaign), ShouldNotBeNil)
		})

		Convey("neither is empty but did'nt match", func() {
			contextCat := []openrtb.ContentCategory{
				openrtb.ContentCategory_IAB1,
				openrtb.ContentCategory_IAB1S5,
			}
			campCat := []openrtb.ContentCategory{
				openrtb.ContentCategory_IAB2,
				openrtb.ContentCategory_IAB3S5,
				openrtb.ContentCategory_IAB6,
			}
			campaign.EXPECT().Category().Return(campCat).AnyTimes()
			context.EXPECT().Category().Return(contextCat).AnyTimes()
			So(catSt.Check(context, campaign), ShouldNotBeNil)
		})

		Convey("neither is empty but they match", func() {
			contextCat := []openrtb.ContentCategory{
				openrtb.ContentCategory_IAB1,
				openrtb.ContentCategory_IAB1S5,
			}
			campCat := []openrtb.ContentCategory{
				openrtb.ContentCategory_IAB1,
				openrtb.ContentCategory_IAB6,
			}
			campaign.EXPECT().Category().Return(campCat).AnyTimes()
			context.EXPECT().Category().Return(contextCat).AnyTimes()
			So(catSt.Check(context, campaign), ShouldBeNil)
		})

	})
}
