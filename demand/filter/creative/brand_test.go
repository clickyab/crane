package creative

import (
	"testing"

	"clickyab.com/crane/demand/entity/mock_entity"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAppBrand_Check(t *testing.T) {
	Convey("checking app brand filter", t, func() {
		ct := gomock.NewController(t)
		brandSt := AppBrand{}
		context := mock_entity.NewMockContext(ct)
		creative := mock_entity.NewMockCreative(ct)
		campaign := mock_entity.NewMockCampaign(ct)
		creative.EXPECT().Campaign().Return(campaign).AnyTimes()

		Convey("campaign have no brands", func() {
			contexBrands := "some brand"
			context.EXPECT().Brand().Return(contexBrands).Times(1)
			campaign.EXPECT().AppBrands().Return([]string{}).Times(1)
			So(brandSt.Check(context, creative), ShouldBeNil)

		})

		Convey("campaign brand didn't match the context brand", func() {
			contexBrands := "some brand"
			context.EXPECT().Brand().Return(contexBrands).Times(1)
			campaign.EXPECT().AppBrands().Return([]string{"some other brand"}).Times(1)
			So(brandSt.Check(context, creative), ShouldNotBeNil)
		})

	})
}
