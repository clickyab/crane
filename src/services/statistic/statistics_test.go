package statistic_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"services/statistic"
	"services/statistic/mock"
	"time"
)

func TestStatisticStore(t *testing.T) {
	statistic.Register(mock.NewMockStatistic)
	Convey("test statistic store", t, func() {
		store := statistic.GetStatisticStore("test_key", 1*time.Hour)
		So(store.Key(), ShouldEqual, "test_key")
		Convey("check inc and dec", func() {

		})
	})
}
