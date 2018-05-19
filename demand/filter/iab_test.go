package filter_test

import (
	"testing"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/entity/mock_entity"
	"clickyab.com/crane/demand/filter"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCategory_Check(t *testing.T) {
	Convey("test iab category filter", t, func() {
		ct := gomock.NewController(t)
		catSt := filter.Category{}
		context := mock_entity.NewMockContext(ct)
		creative := mock_entity.NewMockCreative(ct)
		campaign := mock_entity.NewMockCampaign(ct)
		creative.EXPECT().Campaign().Return(campaign).AnyTimes()

		Convey("campaign category is empty", func() {
			campaign.EXPECT().Category().Return(make([]entity.Category, 0)).AnyTimes()
			So(catSt.Check(context, creative), ShouldBeNil)
		})

		Convey("campaign category is not empty but request cat is empty", func() {
			campCat := makeCat("1", "1-5")
			campaign.EXPECT().Category().Return(campCat).AnyTimes()
			context.EXPECT().Category().Return(make([]entity.Category, 0)).AnyTimes()
			So(catSt.Check(context, creative), ShouldNotBeNil)
		})

		Convey("neither is empty but did'nt match", func() {
			contextCat := makeCat("1", "1-5")
			campCat := makeCat("2", "3-5", "6")
			campaign.EXPECT().Category().Return(campCat).AnyTimes()
			context.EXPECT().Category().Return(contextCat).AnyTimes()
			So(catSt.Check(context, creative), ShouldNotBeNil)
		})

		Convey("neither is empty but they match", func() {
			contextCat := makeCat("1", "1-5")
			campCat := makeCat("2", "3-5", "6", "1-5")
			campaign.EXPECT().Category().Return(campCat).AnyTimes()
			context.EXPECT().Category().Return(contextCat).AnyTimes()
			So(catSt.Check(context, creative), ShouldBeNil)
		})

	})
}

func makeCat(cat ...string) []entity.Category {
	var res = make([]entity.Category, 0)
	for i := range cat {
		res = append(res, entity.Category(cat[i]))
	}
	return res
}
