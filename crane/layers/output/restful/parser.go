package restful

import (
	"net/url"

	"encoding/json"

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

func parse(r *render, imp entity.Impression) error {
	slots := imp.Slots()
	ads := []restAd{}

	for i := range slots {
		slot := slots[i]
		ad := restAd{
			RID: slot.TrackID(),
			// TODO: not sure
			RMaxCPM:      int64(slot.SlotCTR() * float64(slot.WinnerAdvertise().WinnerBID()) * 10),
			RWidth:       slot.Width(),
			RHeight:      slot.Height(),
			RURL:         slot.ShowURL(),
			RLanding:     fetchLanding(slot.WinnerAdvertise().TargetURL()),
			RSlotTrackID: slot.TrackID(),
		}
		ads = append(ads, ad)
	}
	coded, err := json.Marshal(ads)
	if err != nil {
		return err
	}

	r.data = coded
	return nil
}

func fetchLanding(rawURL string) string {
	u, err := url.Parse(rawURL)
	assert.Nil(err)
	return u.Host
}
