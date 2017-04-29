package filter

import (
	"crane/entity"
	"math/rand"
	"testing"
	"time"

	"crane/entity/mock_entity"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFilter(t *testing.T) {
	rand.Seed(time.Now().Unix())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Convey("filters test", t, func() {

		Convey("blacklist filter test", func() {
			Convey("Publisher blackList filter", func() {

				Convey("if it was in black list", func() {
					i := rand.Int63()
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockPublisher := mock_entity.NewMockPublisher(ctrl)
					mockPublisher.EXPECT().ID().Return(i)
					mockImpression.EXPECT().Source().Return(mockPublisher)
					mockClickyabAd.EXPECT().BlackListPublisher().Return([]int64{i, 3, 250})
					So(PublisherBlackList(mockImpression, mockClickyabAd), ShouldBeFalse)
				})

				Convey("if it wasn't in black list", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockPublisher := mock_entity.NewMockPublisher(ctrl)
					mockPublisher.EXPECT().ID().Return(int64(10))
					mockImpression.EXPECT().Source().Return(mockPublisher)
					mockClickyabAd.EXPECT().BlackListPublisher().Return([]int64{12, 20})
					So(PublisherBlackList(mockImpression, mockClickyabAd), ShouldBeTrue)
				})

				Convey("if ad's blacklist is empty", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockClickyabAd.EXPECT().BlackListPublisher().Return([]int64{})
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
			Convey("Publisher whitelist filter", func() {

				Convey("if it was in white list ", func() {
					i := rand.Int63()
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockPublisher := mock_entity.NewMockPublisher(ctrl)
					mockPublisher.EXPECT().ID().Return(i)
					mockImpression.EXPECT().Source().Return(mockPublisher)
					mockClickyabAd.EXPECT().WhiteListPublisher().Return([]int64{i, 10, 100})
					So(PublisherWhiteList(mockImpression, mockClickyabAd), ShouldBeTrue)
				})

				Convey("if it wasn't in white list ", func() {
					i := rand.Int63()
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockPublisher := mock_entity.NewMockPublisher(ctrl)
					mockPublisher.EXPECT().ID().Return(i)
					mockImpression.EXPECT().Source().Return(mockPublisher)
					mockClickyabAd.EXPECT().WhiteListPublisher().Return([]int64{i})
					So(PublisherWhiteList(mockImpression, mockClickyabAd), ShouldBeTrue)
				})

				Convey("if ad's blacklist is empty", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockClickyabAd.EXPECT().WhiteListPublisher().Return([]int64{})
					So(PublisherWhiteList(nil, mockClickyabAd), ShouldBeTrue)
				})

			})

			Convey("Web Category White List ", func() {

				Convey("if it wasnt web type ", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockClickyabAd.EXPECT().Category().Return([]entity.Category{"ads"})
					mockImpression.EXPECT().Category().Return([]entity.Category{"sport", "news"})
					So(Category(mockImpression, mockClickyabAd), ShouldBeFalse)
				})

				Convey("if it was in whitelist", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockClickyabAd.EXPECT().Category().Return([]entity.Category{"sport"})
					mockImpression.EXPECT().Category().Return([]entity.Category{"sport", "news"})
					So(Category(mockImpression, mockClickyabAd), ShouldBeTrue)
				})

				Convey("if it didnt match all whitelist", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockClickyabAd.EXPECT().Category().Return([]entity.Category{"weather", "news"})
					mockImpression.EXPECT().Category().Return([]entity.Category{"sport", "news"})
					So(Category(mockImpression, mockClickyabAd), ShouldBeTrue)
				})

				Convey("if impression cat is empty whitelist", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockClickyabAd.EXPECT().Category().Return([]entity.Category{"weather", "news"})
					mockImpression.EXPECT().Category().Return([]entity.Category{})
					So(Category(mockImpression, mockClickyabAd), ShouldBeFalse)
				})
			})

			Convey("Size white list", func() {

				Convey("if is not in type", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockPublisher := mock_entity.NewMockPublisher(ctrl)
					mockPublisher.EXPECT().AcceptedTarget().Return(entity.TargetApp)
					mockImpression.EXPECT().Source().Return(mockPublisher)
					So(createSizeFilter(entity.TargetVast)(mockImpression, mockClickyabAd), ShouldBeTrue)
				})

				Convey("if is type & in size", func() {
					i := rand.Int()
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockSlot := mock_entity.NewMockSlot(ctrl)
					mockPublisher := mock_entity.NewMockPublisher(ctrl)
					mockSlot.EXPECT().AllowedSize().Return([]int{i, i + 1, i + 2})
					mockClickyabAd.EXPECT().Size().Return(i)
					mockImpression.EXPECT().Slots().Return([]entity.Slot{mockSlot})
					mockPublisher.EXPECT().AcceptedTarget().Return(entity.TargetVast)
					mockImpression.EXPECT().Source().Return(mockPublisher)
					So(createSizeFilter(entity.TargetVast)(mockImpression, mockClickyabAd), ShouldBeTrue)
				})

				Convey("if is not in size", func() {
					i := rand.Int()
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockSlot := mock_entity.NewMockSlot(ctrl)
					mockPublisher := mock_entity.NewMockPublisher(ctrl)
					mockSlot.EXPECT().AllowedSize().Return([]int{i, i + 1, i + 2})
					mockClickyabAd.EXPECT().Size().Return(i + i + 1)
					mockImpression.EXPECT().Slots().Return([]entity.Slot{mockSlot})
					mockPublisher.EXPECT().AcceptedTarget().Return(entity.TargetVast)
					mockImpression.EXPECT().Source().Return(mockPublisher)
					So(createSizeFilter(entity.TargetVast)(mockImpression, mockClickyabAd), ShouldBeFalse)
				})

				Convey("if slot allowed size is empty", func() {
					i := rand.Int()
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockSlot := mock_entity.NewMockSlot(ctrl)
					mockPublisher := mock_entity.NewMockPublisher(ctrl)
					mockSlot.EXPECT().AllowedSize().Return([]int{})
					mockClickyabAd.EXPECT().Size().Return(i + i + 1)
					mockImpression.EXPECT().Slots().Return([]entity.Slot{mockSlot})
					mockPublisher.EXPECT().AcceptedTarget().Return(entity.TargetVast)
					mockImpression.EXPECT().Source().Return(mockPublisher)
					So(createSizeFilter(entity.TargetVast)(mockImpression, mockClickyabAd), ShouldBeFalse)
				})

			})

			Convey("Country white list", func() {

				Convey("if country empty", func() {
					mockClickyabAd := mock_entity.NewMockAdvertise(ctrl)
					mockImpression := mock_entity.NewMockImpression(ctrl)
					mockClickyabAd.EXPECT().Country().Return([]int64{})
					So(Country(mockImpression, mockClickyabAd), ShouldBeTrue)
				})

				Convey("if country ID matched", func() {
					sampleCountry := entity.Country{
						Valid: true,
						ID:    int64(2),
						Name:  "ali",
						ISO:   "iran",
					}
					clickadmock := mock_entity.NewMockAdvertise(ctrl)
					clickadmock.EXPECT().Country().Return([]int64{2})

					locationmock := mock_entity.NewMockLocation(ctrl)
					locationmock.EXPECT().Country().Return(sampleCountry)

					imprmock := mock_entity.NewMockImpression(ctrl)
					imprmock.EXPECT().Location().Return(locationmock)

					So(Country(imprmock, clickadmock), ShouldBeTrue)
				})

				Convey("if country isn't valid", func() {
					sampleCountry := entity.Country{
						Valid: false,
						ID:    int64(2),
						Name:  "ali",
						ISO:   "iran",
					}
					clickadmock := mock_entity.NewMockAdvertise(ctrl)
					clickadmock.EXPECT().Country().Return([]int64{3})

					locationmock := mock_entity.NewMockLocation(ctrl)
					locationmock.EXPECT().Country().Return(sampleCountry)

					imprmock := mock_entity.NewMockImpression(ctrl)
					imprmock.EXPECT().Location().Return(locationmock)
					So(Country(imprmock, clickadmock), ShouldBeFalse)
				})
			})

			Convey("Province white list", func() {

				Convey("if Province ID doesn't matched", func() {
					sampleProvince := entity.Province{
						Valid: true,
						ID:    int64(3),
						Name:  "ali",
					}
					clickadmock := mock_entity.NewMockAdvertise(ctrl)
					clickadmock.EXPECT().Province().Return([]int64{2})

					locationmock := mock_entity.NewMockLocation(ctrl)
					locationmock.EXPECT().Province().Return(sampleProvince)

					imprmock := mock_entity.NewMockImpression(ctrl)
					imprmock.EXPECT().Location().Return(locationmock)
					So(Province(imprmock, clickadmock), ShouldBeFalse)
				})

				Convey("if Province ID matched", func() {
					sampleProvince := entity.Province{
						Valid: true,
						ID:    int64(2),
						Name:  "ali",
					}
					clickadmock := mock_entity.NewMockAdvertise(ctrl)
					clickadmock.EXPECT().Province().Return([]int64{2})

					locationmock := mock_entity.NewMockLocation(ctrl)
					locationmock.EXPECT().Province().Return(sampleProvince)

					imprmock := mock_entity.NewMockImpression(ctrl)
					imprmock.EXPECT().Location().Return(locationmock)
					So(Province(imprmock, clickadmock), ShouldBeTrue)
				})

				Convey("if Province isn't valid", func() {
					sampleProvince := entity.Province{
						Valid: false,
						ID:    int64(2),
						Name:  "ali",
					}
					clickadmock := mock_entity.NewMockAdvertise(ctrl)
					clickadmock.EXPECT().Province().Return([]int64{3})

					locationmock := mock_entity.NewMockLocation(ctrl)
					locationmock.EXPECT().Province().Return(sampleProvince)

					imprmock := mock_entity.NewMockImpression(ctrl)
					imprmock.EXPECT().Location().Return(locationmock)
					So(Province(imprmock, clickadmock), ShouldBeFalse)
				})

				Convey("if Province is empty", func() {
					clickadmock := mock_entity.NewMockAdvertise(ctrl)
					clickadmock.EXPECT().Province().Return([]int64{})
					So(Province(nil, clickadmock), ShouldBeTrue)
				})
			})

			Convey("OS white list", func() {
				osMock := entity.OS{
					ID:     1,
					Mobile: false,
					Valid:  true,
					Name:   "android",
				}

				Convey("if is in white list", func() {
					click := mock_entity.NewMockAdvertise(ctrl)
					click.EXPECT().AllowedOS().Return([]int64{1, 2, 3})

					imp := mock_entity.NewMockImpression(ctrl)
					imp.EXPECT().OS().Return(osMock)

					So(OS(imp, click), ShouldBeTrue)

				})
				Convey("if is not in white list", func() {
					click := mock_entity.NewMockAdvertise(ctrl)
					click.EXPECT().AllowedOS().Return([]int64{2, 3})

					imp := mock_entity.NewMockImpression(ctrl)
					imp.EXPECT().OS().Return(osMock)
					So(OS(imp, click), ShouldBeFalse)

				})

				Convey("if allowed os type is empty", func() {
					click := mock_entity.NewMockAdvertise(ctrl)
					click.EXPECT().AllowedOS().Return([]int64{})
					So(OS(nil, click), ShouldBeTrue)

				})
			})
		})

	})

}
