package restful

import (
	"context"
	"entity"
)

type restAd struct {
	RID     string        `json:"id"`
	RType   entity.AdType `json:"type"`
	RCPM    int64         `json:"cpm"`
	RWidth  int           `json:"width"`
	RHeight int           `json:"height"`
	RCode   string        `json:"code"`

	demand entity.Demand
}

func (ra *restAd) Win() {
	ra.demand.Win(context.Background(), ra.RID, ra.RCPM)
}

func (ra *restAd) Width() int {
	return ra.RWidth
}

func (ra *restAd) Height() int {
	return ra.RHeight
}

func (ra *restAd) ID() string {
	return ra.RID
}

func (ra *restAd) Type() entity.AdType {
	return ra.RType
}

func (ra *restAd) CPM() int64 {
	return ra.RCPM
}

func (ra *restAd) Code() string {
	return ra.RCode
}
