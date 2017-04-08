package restful

import (
	"entity"
	"services/random"
)

type restAd struct {
	RID     string        `json:"id"`
	RType   entity.AdType `json:"type"`
	RMaxCPM int64         `json:"max_cpm"`
	RWidth  int           `json:"width"`
	RHeight int           `json:"height"`
	RURL    string        `json:"url"`

	demand    entity.Demand
	trackID   string
	winnerCPM int64
}

func (ra restAd) URL() string {
	return ra.RURL
}

func (ra *restAd) SetWinnerCPM(w int64) {
	ra.winnerCPM = w
}

func (ra restAd) WinnerCPM() int64 {
	return ra.winnerCPM
}

func (ra *restAd) TrackID() string {
	if ra.trackID == "" {
		ra.trackID = <-random.ID
	}
	return ra.trackID
}

func (ra restAd) Demand() entity.Demand {
	return ra.demand
}

func (ra restAd) Width() int {
	return ra.RWidth
}

func (ra restAd) Height() int {
	return ra.RHeight
}

func (ra restAd) ID() string {
	return ra.RID
}

func (ra restAd) Type() entity.AdType {
	return ra.RType
}

func (ra restAd) MaxCPM() int64 {
	return ra.RMaxCPM
}
