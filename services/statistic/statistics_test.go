package statistic_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"time"

	"clickyab.com/exchange/services/statistic"
	"clickyab.com/exchange/services/statistic/mock"
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
