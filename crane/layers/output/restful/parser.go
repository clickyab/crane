package restful

import (
	"net/url"

	"encoding/json"

	"fmt"

	"time"

	"clickyab.com/crane/crane/entity"
	"github.com/Sirupsen/logrus"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/eav"
)

const (
	restSingleAdEavKey string = `RSA`
	templateKey               = `template`
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

func parse(r *render, imp entity.Impression, cp entity.ClickProvider) error {
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
		go func() {
			renderedAd, err := makeSingleAdData(slot.WinnerAdvertise(), imp, slot, cp)
			if err != nil {
				logrus.Debug("couldn't render ad")
				return
			}
			key := fmt.Sprintf("%s_%s_%s", restSingleAdEavKey, imp.TrackID(), slot.TrackID())
			eav.NewEavStore(key).SetSubKey(templateKey, renderedAd).Save(time.Hour * 24)
		}()
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
