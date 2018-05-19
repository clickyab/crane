package filter_test

import (
	"testing"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/entity/mock_entity"
	"clickyab.com/crane/demand/filter"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProvince_Check(t *testing.T) {
	Convey("checking os filter", t, func() {
		ct := gomock.NewController(t)
		provinceSt := filter.Province{}
		context := mock_entity.NewMockContext(ct)
		creative := mock_entity.NewMockCreative(ct)
		campaign := mock_entity.NewMockCampaign(ct)
		location := mock_entity.NewMockLocation(ct)
		context.EXPECT().Location().Return(location).AnyTimes()
		creative.EXPECT().Campaign().Return(campaign).AnyTimes()

		Convey("campaign province is empty", func() {
			campaign.EXPECT().Province().Return([]string{}).Times(1)
			So(provinceSt.Check(context, creative), ShouldBeNil)
		})

		Convey("campaign not empty and context province not valid", func() {
			campaign.EXPECT().Province().Return([]string{"1", "2"}).Times(1)
			contextProvince := entity.Province{Valid: false}
			location.EXPECT().Province().Return(contextProvince).Times(1)
			So(provinceSt.Check(context, creative), ShouldNotBeNil)
		})

		Convey("context location is iran and campaign does not have iran id (1)", func() {
			campaign.EXPECT().Province().Return([]string{"3", "2"}).Times(3)
			contextCountry := entity.Country{Valid: true, ISO: "IR"}
			contextProvince := entity.Province{Valid: true}
			location.EXPECT().Country().Return(contextCountry).Times(1)
			location.EXPECT().Province().Return(contextProvince).Times(2)
			So(provinceSt.Check(context, creative), ShouldNotBeNil)
		})

		Convey("context location is iran and campaign match iran id (1)", func() {
			campaign.EXPECT().Province().Return([]string{"1", "2"}).Times(2)
			contextCountry := entity.Country{Valid: true, ISO: "IR"}
			contextProvince := entity.Province{Valid: true}
			location.EXPECT().Country().Return(contextCountry).Times(1)
			location.EXPECT().Province().Return(contextProvince).Times(2)
			So(provinceSt.Check(context, creative), ShouldBeNil)
		})

		Convey("context location is not iran and campaign", func() {
			campaign.EXPECT().Province().Return([]string{"1", "2"}).Times(2)
			contextCountry := entity.Country{Valid: true, ISO: "AF"}
			contextProvince := entity.Province{Valid: true, ID: 2}
			location.EXPECT().Country().Return(contextCountry).Times(1)
			location.EXPECT().Province().Return(contextProvince).Times(2)
			So(provinceSt.Check(context, creative), ShouldBeNil)
		})

	})
}
