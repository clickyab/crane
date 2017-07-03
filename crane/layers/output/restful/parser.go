package restful

import (
	"net/url"

	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/assert"
)

type restAd struct {
	RID          string `json:"id"`
	RMaxCPM      int64  `json:"max_cpm"`
	RWidth       int    `json:"width"`
	RHeight      int    `json:"height"`
	RURL         string `json:"url"`
	RLanding     string `json:"landing"`
	RSlotTrackID string `json:"slot_track_id"`
}

func parser(imp entity.Impression) []restAd {
	slots := imp.Slots()
	ads := []restAd{}

	for i := range slots {
		slot := slots[i]
		ad := restAd{
			RID:          slot.TrackID(),
			RMaxCPM:      int64(slot.SlotCTR() * 10 * float64(slot.WinnerAdvertise().WinnerBID())),
			RWidth:       slot.Width(),
			RHeight:      slot.Height(),
			RURL:         slot.ShowURL(),
			RLanding:     fetchLanding(fetchLanding(slot.WinnerAdvertise().TargetURL())),
			RSlotTrackID: slot.TrackID(),
		}
		ads = append(ads, ad)
	}
	return ads
}

func fetchLanding(rawURL string) string {
	u, err := url.Parse(rawURL)
	assert.Nil(err)
	return u.Host
}
