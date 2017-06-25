package filter

import (
	"math/rand"
	"testing"
	"time"

	"clickyab.com/crane/crane/entity"

	"clickyab.com/crane/crane/entity/mock_entity"

	"github.com/clickyab/services/random"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFilter(t *testing.T) {
	rand.Seed(time.Now().Unix())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Convey("filters test", t, func() {
		clickyabAd := mock_entity.NewMockAdvertise(ctrl)
		impression := mock_entity.NewMockImpression(ctrl)
		campaign := mock_entity.NewMockCampaign(ctrl)
		clickyabAd.EXPECT().Campaign().Return(campaign).AnyTimes()
		Convey("blacklist filter test", func() {
			Convey("Publisher blackList filter", func() {

				Convey("if it was in black list", func() {
					i := <-random.ID
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockPublisher := mock_entity.NewMockPublisher(ctrl)
					mockCampaign := mock_entity.NewMockCampaign(ctrl)
					mockCampaign.EXPECT().BlackListPublisher().Return([]string{i, "one", "two"})
					mockPublisher.EXPECT().Name().Return(i)
					mockImpression.EXPECT().Source().Return(mockPublisher)
					mockClickyabAd.EXPECT().Campaign().Return(mockCampaign)

					So(PublisherBlackList(mockImpression, mockClickyabAd), ShouldBeFalse)
				})

				Convey("if it wasn't in black list", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockPublisher := mock_entity.NewMockPublisher(ctrl)
					mockCampaign := mock_entity.NewMockCampaign(ctrl)
					mockCampaign.EXPECT().BlackListPublisher().Return([]string{"one", "two"})
					mockClickyabAd.EXPECT().Campaign().Return(mockCampaign)
					mockPublisher.EXPECT().Name().Return(<-random.ID)
					mockImpression.EXPECT().Source().Return(mockPublisher)
					So(PublisherBlackList(mockImpression, mockClickyabAd), ShouldBeTrue)
				})

				Convey("if ad's blacklist is empty", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockCampaign := mock_entity.NewMockCampaign(ctrl)
					mockCampaign.EXPECT().BlackListPublisher().Return([]string{})
					mockClickyabAd.EXPECT().Campaign().Return(mockCampaign)
					So(PublisherBlackList(nil, mockClickyabAd), ShouldBeTrue)
				})

			})
			Convey("Target testing", func() {

				Convey("impressions source type is app and matches ad type", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockPublisher := mock_entity.NewMockPublisher(ctrl)
					mockCampaign := mock_entity.NewMockCampaign(ctrl)
					mockCampaign.EXPECT().Target().Return([]entity.Target{entity.TargetApp, entity.TargetVast})
					mockPublisher.EXPECT().AcceptedTarget().Return(entity.TargetApp)
					mockClickyabAd.EXPECT().Campaign().Return(mockCampaign)
					mockImpression.EXPECT().Source().Return(mockPublisher)
					So(Target(mockImpression, mockClickyabAd), ShouldBeTrue)
				})

				Convey("impressions source type is web and matches ad type", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockPublisher := mock_entity.NewMockPublisher(ctrl)
					mockCampaign := mock_entity.NewMockCampaign(ctrl)
					mockCampaign.EXPECT().Target().Return([]entity.Target{entity.TargetWeb, entity.TargetVast})
					mockPublisher.EXPECT().AcceptedTarget().Return(entity.TargetWeb)
					mockClickyabAd.EXPECT().Campaign().Return(mockCampaign)
					mockImpression.EXPECT().Source().Return(mockPublisher)
					So(Target(mockImpression, mockClickyabAd), ShouldBeTrue)
				})

				Convey("impressions source type is vast and matches ad type", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockPublisher := mock_entity.NewMockPublisher(ctrl)
					mockCampaign := mock_entity.NewMockCampaign(ctrl)
					mockCampaign.EXPECT().Target().Return([]entity.Target{entity.TargetVast, entity.TargetApp})
					mockPublisher.EXPECT().AcceptedTarget().Return(entity.TargetVast)
					mockClickyabAd.EXPECT().Campaign().Return(mockCampaign)
					mockImpression.EXPECT().Source().Return(mockPublisher)
					So(Target(mockImpression, mockClickyabAd), ShouldBeTrue)
				})

				Convey("impressions source type is app and doesn't match ad type", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockPublisher := mock_entity.NewMockPublisher(ctrl)
					mockCampaign := mock_entity.NewMockCampaign(ctrl)
					mockCampaign.EXPECT().Target().Return([]entity.Target{entity.TargetWeb, entity.TargetVast})
					mockPublisher.EXPECT().AcceptedTarget().Return(entity.TargetApp)
					mockClickyabAd.EXPECT().Campaign().Return(mockCampaign)
					mockImpression.EXPECT().Source().Return(mockPublisher)
					So(Target(mockImpression, mockClickyabAd), ShouldBeFalse)
				})

				Convey("impressions source type is web and doesn't match ad type", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockPublisher := mock_entity.NewMockPublisher(ctrl)
					mockCampaign := mock_entity.NewMockCampaign(ctrl)
					mockCampaign.EXPECT().Target().Return([]entity.Target{entity.TargetVast})
					mockPublisher.EXPECT().AcceptedTarget().Return(entity.TargetWeb)
					mockClickyabAd.EXPECT().Campaign().Return(mockCampaign)
					mockImpression.EXPECT().Source().Return(mockPublisher)
					So(Target(mockImpression, mockClickyabAd), ShouldBeFalse)
				})

				Convey("impressions source type is vast and doesn't match ad type", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockPublisher := mock_entity.NewMockPublisher(ctrl)
					mockCampaign := mock_entity.NewMockCampaign(ctrl)
					mockCampaign.EXPECT().Target().Return([]entity.Target{entity.TargetWeb, entity.TargetApp})
					mockPublisher.EXPECT().AcceptedTarget().Return(entity.TargetVast)
					mockClickyabAd.EXPECT().Campaign().Return(mockCampaign)
					mockImpression.EXPECT().Source().Return(mockPublisher)
					So(Target(mockImpression, mockClickyabAd), ShouldBeFalse)
				})
			})
		})

		Convey("whitelist filter test", func() {
			//Convey("Publisher whitelist filter", func() {
			//
			//	Convey("if it was in white list ", func() {
			//		id := <-random.ID
			//		mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
			//		mockImpression := mock_entity.NewMockImpression(ctrl)
			//		mockPublisher := mock_entity.NewMockPublisher(ctrl)
			//		mockCampaign := mock_entity.NewMockCampaign(ctrl)
			//		mockCampaign.EXPECT().WhiteListPublisher().Return([]string{id, "one", "two"})
			//		mockPublisher.EXPECT().Name().Return(id)
			//		mockImpression.EXPECT().Source().Return(mockPublisher)
			//		So(PublisherWhiteList(mockImpression, mockClickyabAd), ShouldBeTrue)
			//	})
			//
			//	Convey("if it wasn't in white list ", func() {
			//		id := <-random.ID
			//		mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
			//		mockImpression := mock_entity.NewMockImpression(ctrl)
			//		mockPublisher := mock_entity.NewMockPublisher(ctrl)
			//		mockCampaign := mock_entity.NewMockCampaign(ctrl)
			//		mockCampaign.EXPECT().WhiteListPublisher().Return([]string{})
			//		mockClickyabAd.EXPECT().Campaign().Return(mockCampaign)
			//		mockPublisher.EXPECT().Name().Return(id)
			//		mockImpression.EXPECT().Source().Return(mockPublisher)
			//		So(PublisherWhiteList(mockImpression, mockClickyabAd), ShouldBeTrue)
			//	})
			//
			//	Convey("if ad's blacklist is empty", func() {
			//
			//		mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
			//		mockCampaign := mock_entity.NewMockCampaign(ctrl)
			//		mockClickyabAd.EXPECT().Campaign().Return(mockCampaign)
			//		mockCampaign.EXPECT().WhiteListPublisher().Return([]string{})
			//		So(PublisherWhiteList(nil, mockClickyabAd), ShouldBeTrue)
			//	})
			//
			//})

			Convey("Web Category White List ", func() {

				Convey("if it wasnt web type ", func() {

					campaign.EXPECT().Category().Return([]entity.Category{"ads"})
					impression.EXPECT().Category().Return([]entity.Category{"sport", "news"})
					So(Category(impression, clickyabAd), ShouldBeFalse)
				})

				Convey("if it was in whitelist", func() {
					campaign.EXPECT().Category().Return([]entity.Category{"sport"})
					impression.EXPECT().Category().Return([]entity.Category{"sport", "news"})
					So(Category(impression, clickyabAd), ShouldBeTrue)
				})

				Convey("if it didnt match all whitelist", func() {
					campaign.EXPECT().Category().Return([]entity.Category{"weather", "news"})
					impression.EXPECT().Category().Return([]entity.Category{"sport", "news"})
					So(Category(impression, clickyabAd), ShouldBeTrue)
				})

				Convey("if impression cat is empty whitelist", func() {
					campaign.EXPECT().Category().Return([]entity.Category{"weather", "news"})
					impression.EXPECT().Category().Return([]entity.Category{})
					So(Category(impression, clickyabAd), ShouldBeFalse)
				})
			})

			Convey("Country white list", func() {
				Convey("if country empty", func() {
					campaign.EXPECT().Country().Return([]string{})
					So(Country(impression, clickyabAd), ShouldBeTrue)
				})

				Convey("if country ID matched", func() {
					sampleCountry := entity.Country{
						Valid: true,
						Name:  "ali",
						ISO:   "iran",
					}
					campaign.EXPECT().Country().Return([]string{"ali"})

					location := mock_entity.NewMockLocation(ctrl)
					location.EXPECT().Country().Return(sampleCountry)

					imprmock := mock_entity.NewMockImpression(ctrl)
					imprmock.EXPECT().Location().Return(location)

					So(Country(imprmock, clickyabAd), ShouldBeTrue)
				})

				Convey("if country isn't valid", func() {
					sampleCountry := entity.Country{
						Valid: false,
						Name:  "ali",
						ISO:   "iran",
					}
					campaign.EXPECT().Country().Return([]string{""})

					locationmock := mock_entity.NewMockLocation(ctrl)
					locationmock.EXPECT().Country().Return(sampleCountry)

					imprmock := mock_entity.NewMockImpression(ctrl)
					imprmock.EXPECT().Location().Return(locationmock)
					So(Country(imprmock, clickyabAd), ShouldBeFalse)
				})
			})

			Convey("Province white list", func() {

				Convey("if Province ID doesn't matched", func() {
					sampleProvince := entity.Province{
						Valid: true,
						Name:  "ali",
					}
					campaign.EXPECT().Province().Return([]string{""})

					locationmock := mock_entity.NewMockLocation(ctrl)
					locationmock.EXPECT().Province().Return(sampleProvince)

					imprmock := mock_entity.NewMockImpression(ctrl)
					imprmock.EXPECT().Location().Return(locationmock)
					So(Province(imprmock, clickyabAd), ShouldBeFalse)
				})

				Convey("if Province ID matched", func() {
					sampleProvince := entity.Province{
						Valid: true,
						Name:  "ali",
					}
					campaign.EXPECT().Province().Return([]string{"ali"})

					locationmock := mock_entity.NewMockLocation(ctrl)
					locationmock.EXPECT().Province().Return(sampleProvince)

					imprmock := mock_entity.NewMockImpression(ctrl)
					imprmock.EXPECT().Location().Return(locationmock)
					So(Province(imprmock, clickyabAd), ShouldBeTrue)

				})

				Convey("if Province isn't valid", func() {
					sampleProvince := entity.Province{
						Valid: false,
						Name:  "ali",
					}
					campaign.EXPECT().Province().Return([]string{""})

					locationmock := mock_entity.NewMockLocation(ctrl)
					locationmock.EXPECT().Province().Return(sampleProvince)

					imprmock := mock_entity.NewMockImpression(ctrl)
					imprmock.EXPECT().Location().Return(locationmock)
					So(Province(imprmock, clickyabAd), ShouldBeFalse)
				})

				Convey("if Province is empty", func() {
					campaign.EXPECT().Province().Return([]string{})
					So(Province(nil, clickyabAd), ShouldBeTrue)
				})
			})

			Convey("OS white list", func() {
				osMock := entity.OS{
					Mobile: false,
					Valid:  true,
					Name:   "android",
				}

				Convey("if is in white list", func() {
					campaign.EXPECT().AllowedOS().Return([]string{})

					imp := mock_entity.NewMockImpression(ctrl)
					imp.EXPECT().OS().Return(osMock).AnyTimes()

					So(OS(imp, clickyabAd), ShouldBeTrue)

				})
				Convey("if is not in white list", func() {
					campaign.EXPECT().AllowedOS().Return([]string{"ios"})

					imp := mock_entity.NewMockImpression(ctrl)
					imp.EXPECT().OS().Return(osMock).AnyTimes()
					So(OS(imp, clickyabAd), ShouldBeFalse)

				})

				Convey("if allowed os type is empty", func() {
					campaign.EXPECT().AllowedOS().Return([]string{}).AnyTimes()
					So(OS(nil, clickyabAd), ShouldBeTrue)

				})
			})
		})

	})

}
