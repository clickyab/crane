package restful

import (
	"io"
	"net/url"

	"encoding/json"

	"clickyab.com/crane/crane/entity"

	"github.com/clickyab/services/assert"
)

type render struct {
}

// need to register render{} somewhere
func (r render) Render(w io.Writer, imp entity.Impression, cp entity.ClickProvider) error {

	slots := imp.Slots()
	ads := []restAd{}

	for i := range slots {
		slot := slots[i]
		ad := restAd{
			RID: slot.WinnerAdvertise().ID(),
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

	_, err = w.Write(coded)
	return err
}

type restAd struct {
	RID          string `json:"id"`
	RMaxCPM      int64  `json:"max_cpm"`
	RWidth       int    `json:"width"`
	RHeight      int    `json:"height"`
	RURL         string `json:"url"`
	RLanding     string `json:"landing"`
	RSlotTrackID string `json:"slot_track_id"`
}

func fetchLanding(rawURL string) string {
	u, err := url.Parse(rawURL)
	assert.Nil(err)
	return u.Host
}
