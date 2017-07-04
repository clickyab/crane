package local

import (
	"clickyab.com/crane/crane/entity"

	"github.com/clickyab/services/assert"
)

// Slot Slot
type Slot struct {
	FID      string `json:"id"`
	FWidth   int    `json:"width"`
	FHeight  int    `json:"height"`
	FTrackID string `json:"track_id"`
	slotCTR  float64

	attribute map[string]interface{}
	winnerAd  interface{}
	showURL   string
}

// ID is slot's ID
func (rs *Slot) ID() string {
	return rs.FID
}

// TrackID is slot's random trackID
func (rs *Slot) TrackID() string {
	return rs.FTrackID
}

// Width is slot's width
// @required slot's Width
func (rs *Slot) Width() int {
	return rs.FWidth
}

// Height is slot's Height
// @required FHeight
func (rs *Slot) Height() int {
	return rs.FHeight
}

// SlotCTR is slot's ctr
func (rs *Slot) SlotCTR() float64 {
	return rs.slotCTR
}

// SetSlotCTR set's slot's ctr
func (rs *Slot) SetSlotCTR(ctr float64) {
	rs.slotCTR = ctr
}

// SetWinnerAdvertise Sets slot's Winner Advertise
func (rs *Slot) SetWinnerAdvertise(ad entity.Advertise) {
	rs.winnerAd = ad
}

// WinnerAdvertise gets slot's Winner Advertise
func (rs *Slot) WinnerAdvertise() entity.Advertise {
	ad, ok := rs.winnerAd.(entity.Advertise)
	assert.True(ok)
	return ad
}

// SetShowURL set slot's show url
func (rs *Slot) SetShowURL(url string) {
	rs.showURL = url
}

// ShowURL returns showURL, the url' response is and html rendered ad
func (rs *Slot) ShowURL() string {
	return rs.showURL
}

// IsSizeAllowed says if its allowed by owner size
func (rs *Slot) IsSizeAllowed(width, height int) bool {
	w, h := rs.FWidth, rs.FHeight
	if w == width && h == height {
		return true
	}
	return false
}

// Attribute slots attribute
func (rs *Slot) Attribute() map[string]interface{} {
	return rs.attribute
}
