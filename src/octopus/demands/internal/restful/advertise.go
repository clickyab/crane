package restful

import (
	"octopus/exchange"
	"services/random"
)

type restAd struct {
	RID     string `json:"id"`
	RMaxCPM int64  `json:"max_cpm"`
	RWidth  int    `json:"width"`
	RHeight int    `json:"height"`
	RURL    string `json:"url"`

	demand    exchange.Demand
	trackID   string
	winnerCPM int64
	rates     []exchange.Rate
	landing   string
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

func (ra restAd) Demand() exchange.Demand {
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

func (ra restAd) MaxCPM() int64 {
	return ra.RMaxCPM
}

func (ra restAd) Rates() []exchange.Rate {
	return ra.rates
}

func (ra restAd) Landing() string {
	return ra.landing
}
