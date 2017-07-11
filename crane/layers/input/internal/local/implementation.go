package local

import (
	"net"

	"clickyab.com/crane/crane/entity"
	"github.com/Sirupsen/logrus"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
)

var originalUnderFloor = config.RegisterBoolean("crane.input.rest.under_floor", false, "its used when publisher's underfloor isn't set")

const acceptedTarget string = "accepted_target"

// FloorCPM is publisher's cpm floor
// @required FFloorCPM
func (rp *Publisher) FloorCPM() int64 {
	return rp.FFloorCPM
}

// SoftFloorCPM is publisher soft floor cpm
// @required FSoftFloorCPM
func (rp *Publisher) SoftFloorCPM() int64 {
	return rp.FSoftFloorCPM
}

// UnderFloor is publisher underfloor
func (rp *Publisher) UnderFloor() bool {
	if rp.UnderFloor == nil {
		return originalUnderFloor.Bool()
	}
	return *rp.FUnderFloor
}

// Name is publisher name
// @required FName
func (rp *Publisher) Name() string {
	return rp.FName
}

// AcceptedTarget is publisher's target (web, vast, app, native)
func (rp *Publisher) AcceptedTarget() entity.Target {
	t, ok := rp.FAttributes[acceptedTarget].(entity.Target)
	if !ok {
		return entity.TargetInvalid
	}

	return t
}

// Attributes is publisher's Attributes
func (rp *Publisher) Attributes() map[string]interface{} {
	return rp.FAttributes
}

// BIDType is publisher's bid type, rest is cpm
func (rp *Publisher) BIDType() entity.BIDType {
	return entity.BIDTypeCPM
}

// MinCPC is publisher's minimum cpc
func (rp *Publisher) MinCPC() int64 {
	logrus.Panic("rest type shouldn't have minCPC")
	return 0
}

// AcceptedTypes is publisher's accepted types (dyn, banner, video, html, native)
func (rp *Publisher) AcceptedTypes() []entity.AdType {
	at, ok := rp.FAttributes["accepted_types"].([]entity.AdType)
	if !ok {
		return nil
	}

	return at
}

// Supplier is publisher's supplier
// @required FSupplier
func (rp *Publisher) Supplier() string {
	return rp.FSupplier
}

// Country is  country
func (rl Location) Country() entity.Country {
	return rl.FCountry
}

// Province is request's Province
func (rl Location) Province() entity.Province {
	return rl.FProvince
}

// LatLon is request's LatLon
func (rl Location) LatLon() entity.LatLon {
	return rl.FLatLon
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

func (r *request) IP() net.IP {
	return r.ip
}

func (r *request) OS() entity.OS {
	return r.os
}

func (r *request) ClientID() string {
	return r.client
}

func (r *request) Protocol() string {
	return r.protocol
}

func (r *request) UserAgent() string {
	return r.userAgent
}

func (r *request) Location() entity.Location {
	return r.location
}

func (r *request) Attributes() map[string]string {
	return r.attr
}
