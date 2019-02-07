package campaign

import (
	"testing"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/entity/mock_entity"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAreaInGlob_Check(t *testing.T) {
	Convey("test area in glob filter", t, func() {
		ct := gomock.NewController(t)
		globSt := AreaInGlob{}
		context := mock_entity.NewMockContext(ct)
		campaign := mock_entity.NewMockCampaign(ct)
		location := mock_entity.NewMockLocation(ct)
		context.EXPECT().Location().Return(location).AnyTimes()

		Convey("campaign is not regional", func() {
			var a float64
			latLon := entity.LatLon{}
			location.EXPECT().LatLon().Return(latLon).AnyTimes()
			campaign.EXPECT().LatLon().Return(false, a, a, a).AnyTimes()
			So(globSt.Check(context, campaign), ShouldBeNil)
		})

		Convey("campaign is regional but context lat lon not detected", func() {
			latLon := entity.LatLon{Valid: false}
			location.EXPECT().LatLon().Return(latLon).AnyTimes()
			campaign.EXPECT().LatLon().Return(true, 1.5, 1.5, 1.5).AnyTimes()
			So(globSt.Check(context, campaign), ShouldNotBeNil)
		})

		Convey("campaign is regional and context lat lon detected", func() {
			latLon := entity.LatLon{Valid: true, Lat: 35.0, Lon: 36.0}
			location.EXPECT().LatLon().Return(latLon).AnyTimes()
			campaign.EXPECT().LatLon().Return(true, 33.9, 34.9, 1000.0).AnyTimes()
			So(globSt.Check(context, campaign), ShouldBeNil)
		})

		Convey("campaign is regional and context lat lon detected but not match", func() {
			latLon := entity.LatLon{Valid: true, Lat: 35.0, Lon: 36.0}
			location.EXPECT().LatLon().Return(latLon).AnyTimes()
			campaign.EXPECT().LatLon().Return(true, 33.9, 34.9, 100.0).AnyTimes()
			So(globSt.Check(context, campaign), ShouldNotBeNil)
		})
	})
}
