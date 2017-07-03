package vast

import (
	"net"
	"net/http"

	"clickyab.com/crane/crane/entity"
)

func (v *imp) Request() *http.Request {
	return v.FRequest
}

func (v *imp) TrackID() string {
	return v.FTrackID
}

func (v *imp) ClientID() string {
	return v.FClientID
}

func (v *imp) IP() net.IP {
	return v.FIP
}

func (v *imp) UserAgent() string {
	return v.FUserAgenr
}

func (v *imp) Source() entity.Publisher {
	return v.FPublisher
}

func (v *imp) Location() entity.Location {
	return v.FLocation
}

func (v *imp) OS() entity.OS {
	return v.FOS
}

func (v *imp) Slots() []entity.Slot {
	if v.vDum == nil {
		v.vDum = make([]entity.Slot, len(v.FSlots))
		for j := range v.FSlots {
			v.vDum[j] = v.FSlots[j]
		}
	}
	return v.vDum
}

func (v *imp) Category() []entity.Category {
	return v.FCategories
}

func (v *imp) Attributes() map[string]interface{} {
	return v.FAttr
}
