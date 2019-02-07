package campaign

import (
	"testing"

	"clickyab.com/crane/demand/entity/mock_entity"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPublisher_Check(t *testing.T) {
	Convey("test white list and black list publisher", t, func() {
		ct := gomock.NewController(t)
		whiteSt := WhiteList{}
		blackSt := BlackList{}
		context := mock_entity.NewMockContext(ct)
		campaign := mock_entity.NewMockCampaign(ct)
		publisher := mock_entity.NewMockPublisher(ct)
		context.EXPECT().Publisher().Return(publisher).AnyTimes()

		Convey("campaign white list is empty and context publisher id exists", func() {
			campaign.EXPECT().WhiteListPublisher().Return([]string{}).Times(1)
			publisher.EXPECT().ID().Return(int64(1)).Times(1)
			So(whiteSt.Check(context, campaign), ShouldBeNil)
		})

		Convey("campaign white list is not empty and context publisher id exists and they match", func() {
			campaign.EXPECT().WhiteListPublisher().Return([]string{"1", "2"}).Times(1)
			publisher.EXPECT().ID().Return(int64(1)).Times(1)
			So(whiteSt.Check(context, campaign), ShouldBeNil)
		})

		Convey("campaign white list is not empty and context publisher id exists and they not match", func() {
			campaign.EXPECT().WhiteListPublisher().Return([]string{"1", "2"}).Times(1)
			publisher.EXPECT().ID().Return(int64(3)).Times(1)
			So(whiteSt.Check(context, campaign), ShouldNotBeNil)
		})

		Convey("campaign black list is empty and context publisher id exists", func() {
			campaign.EXPECT().BlackListPublisher().Return([]string{}).Times(1)
			publisher.EXPECT().ID().Return(int64(1)).Times(1)
			So(blackSt.Check(context, campaign), ShouldBeNil)
		})

		Convey("campaign black list is not empty and context publisher id exists and they match", func() {
			campaign.EXPECT().BlackListPublisher().Return([]string{"1", "2"}).Times(1)
			publisher.EXPECT().ID().Return(int64(1)).Times(1)
			So(blackSt.Check(context, campaign), ShouldNotBeNil)
		})

		Convey("campaign black list is not empty and context publisher id exists and they not match", func() {
			campaign.EXPECT().BlackListPublisher().Return([]string{"1", "2"}).Times(1)
			publisher.EXPECT().ID().Return(int64(3)).Times(1)
			So(blackSt.Check(context, campaign), ShouldBeNil)
		})
	})
}
