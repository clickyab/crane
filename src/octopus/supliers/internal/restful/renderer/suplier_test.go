package renderer

import (
	"testing"

	"octopus/exchange/mock_exchange"

	"fmt"

	"octopus/exchange"

	"encoding/json"

	"bytes"

	"net/http"

	"github.com/fatih/structs"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

type testResponseWriter struct {
	headers http.Header
	status  int
	buff    bytes.Buffer
}

func (rw *testResponseWriter) Header() http.Header {
	return rw.headers
}

func (rw *testResponseWriter) Write(p []byte) (int, error) {
	return rw.buff.Write(p)
}

func (rw *testResponseWriter) WriteHeader(i int) {
	rw.status = i
}

func TestSupplier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Convey("reder test", t, func() {
		supplier := mock_exchange.NewMockSupplier(ctrl)
		supplier.EXPECT().Share().Return(100).AnyTimes()

		impression := mock_exchange.NewMockImpression(ctrl)

		trackIDs := []string{"aaa", "bbb", "ccc"}
		slots := []exchange.Slot{}
		ads := map[string]exchange.Advertise{}

		// generating slots and ads
		for _, trackID := range trackIDs {
			slot := mock_exchange.NewMockSlot(ctrl)
			ad := mock_exchange.NewMockAdvertise(ctrl)
			demand := mock_exchange.NewMockDemand(ctrl)

			demand.EXPECT().Name().Return("daniel").AnyTimes()

			ad.EXPECT().WinnerCPM().Return(int64(100)).AnyTimes()
			ad.EXPECT().Width().Return(20).AnyTimes()
			ad.EXPECT().Height().Return(15).AnyTimes()
			ad.EXPECT().Landing().Return("clickyab.ir").AnyTimes()
			ad.EXPECT().Demand().Return(demand).AnyTimes()
			ad.EXPECT().TrackID().Return(trackID).AnyTimes()
			ad.EXPECT().URL().Return("www.ad_url.com").AnyTimes()

			fallback := fmt.Sprintf("www.%s.com", trackID)
			slot.EXPECT().Fallback().Return(fallback).AnyTimes()
			slot.EXPECT().TrackID().Return(trackID).AnyTimes()
			slot.EXPECT().Width().Return(20).AnyTimes()
			slot.EXPECT().Height().Return(15).AnyTimes()

			slots = append(slots, slot)
			ads[trackID] = ad
		}

		impression.EXPECT().Slots().Return(slots).AnyTimes()
		impression.EXPECT().Scheme().Return("http").AnyTimes()

		rf := restful{
			pixelPattern: "/pixel/%s/%s",
			sup:          supplier,
		}

		var w = testResponseWriter{
			headers: make(http.Header),
		}

		err := rf.Render(impression, ads, &w)
		So(err, ShouldBeNil)

		resultStruct := []*dumbAd{}
		result := []map[string]interface{}{}
		err = json.Unmarshal(w.buff.Bytes(), &resultStruct)
		So(err, ShouldBeNil)

		for i := range resultStruct {
			resultStruct[i].Code = ""
			result = append(result, structs.New(*resultStruct[i]).Map())
		}

		So(result, ShouldResemble, expected)
		So(w.status, ShouldEqual, http.StatusOK)
	})
}

var expected = []map[string]interface{}{
	{
		"is_filled":   true,
		"landing":     "clickyab.ir",
		"track_id":    "aaa",
		"ad_track_id": "aaa",
		"winner":      int64(0),
		"width":       20,
		"height":      15,
		"code":        "",
	},
	{
		"is_filled":   true,
		"landing":     "clickyab.ir",
		"track_id":    "bbb",
		"ad_track_id": "bbb",
		"winner":      int64(0),
		"width":       20,
		"height":      15,
		"code":        "",
	},
	{
		"is_filled":   true,
		"landing":     "clickyab.ir",
		"track_id":    "ccc",
		"ad_track_id": "ccc",
		"winner":      int64(0),
		"width":       20,
		"height":      15,
		"code":        "",
	},
}
