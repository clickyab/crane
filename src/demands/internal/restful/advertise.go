package restful

import (
	"entity"
	"fmt"
	"net/url"
	"services/random"
)

type restAd struct {
	RID     string `json:"id"`
	RMaxCPM int64  `json:"max_cpm"`
	RWidth  int    `json:"width"`
	RHeight int    `json:"height"`
	RURL    string `json:"url"`

	demand    entity.Demand
	trackID   string
	winnerCPM int64
}

func (ra restAd) Morph(pixel string) entity.DumbAd {
	res := entity.DumbAd{
		MaxCPM: ra.MaxCPM(),
		Height: ra.Height(),
		ID:     ra.ID(),
		Width:  ra.Width(),
	}
	ul := ra.URL()
	u, err := url.Parse(ra.URL())
	if err == nil {
		v := u.Query()
		v.Add("win", fmt.Sprint(ra.winnerCPM))
		u.RawQuery = v.Encode()
		ul = u.String()
	}
	res.Code = fmt.Sprintf(`<img src="%s"><iframe src="%s"></iframe>`, pixel, ul)
	return res
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

func (ra restAd) MaxCPM() int64 {
	return ra.RMaxCPM
}
