package restful

import (
	"net"
	"net/http"

	"clickyab.com/crane/crane/entity"
)

// Request is Request obviously
// this should be filled while decoding
func (ri *restInput) Request() *http.Request {
	return ri.request.r
}

// TrackID is TrackID obviously
// @required FTrackID
func (ri *restInput) TrackID() string {
	return ri.FTrackID
}

// ClientID is ClientID obviously
// @required UserTrackID
func (ri *restInput) ClientID() string {
	return ri.UserTrackID
}

// IP is IP obviously
// @required FIP
func (ri *restInput) IP() net.IP {
	return ri.FIP
}

// UserAgent is UserAgent obviously
// this should be while decoding
func (ri *restInput) UserAgent() string {
	return ri.FUserAgent
}

// Source is Source obviously
// @required FSource
func (ri *restInput) Source() entity.Publisher {
	return ri.FSource
}

// Location is Location obviously
func (ri *restInput) Location() entity.Location {
	return ri.Locationn
}

// OS is OS obviously
func (ri *restInput) OS() entity.OS {
	if ri.os.Valid == false {
		ri.os = entity.OsFromUA(ri.r.UserAgent())
	}
	return ri.os
}

// Attributes is Attributes obviously
func (ri *restInput) Attributes() map[string]interface{} {
	return ri.FAttributes
}

// Slots is Slots obviously
// @required FSlots
func (ri *restInput) Slots() []entity.Slot {
	resp := []entity.Slot{}
	for i := range ri.FSlots {
		resp = append(resp, ri.FSlots[i])
	}
	return resp
}

// Category is Category obviously
func (ri *restInput) Category() []entity.Category {
	resp := []entity.Category{}
	for i := range ri.FCategory {
		resp = append(resp, ri.FCategory[i])
	}
	return resp
}
