package native

import (
	"net"
	"net/http"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/layers/input/internal/local"
	"github.com/clickyab/services/ip2location"
)

func (i *imp) Request() *http.Request {
	return i.FRequest
}

func (i *imp) TrackID() string {
	return i.FTrackID
}

func (i *imp) ClientID() string {
	return i.FClientID
}

func (i *imp) IP() net.IP {
	return i.FIP
}

func (i *imp) UserAgent() string {
	return i.FUA
}

func (i *imp) Source() entity.Publisher {
	return i.FPub
}

func (i *imp) Location() entity.Location {
	return i.FLocation
}

func (i *imp) OS() entity.OS {
	return i.FOS
}

func (i *imp) Slots() []entity.Slot {
	if i.nDum == nil {
		i.nDum = make([]entity.Slot, len(i.FSlots))
		for j := range i.FSlots {
			i.nDum[j] = i.FSlots[j]
		}
	}
	return i.nDum
}

func (i *imp) Category() []entity.Category {
	return i.FCategories
}

func (i *imp) Attributes() map[string]interface{} {
	return i.FAttr
}

func (i *imp) extractData() {
	d := ip2location.IP2Location(i.FIP.String())
	i.FLocation = &local.Location{
		TheCountry: entity.Country{
			Name:  d.CountryLong,
			ISO:   d.CountryShort,
			Valid: d.CountryLong != "-",
		},

		TheProvince: entity.Province{
			Valid: d.Region != "-",
			Name:  d.Region,
		},

		TheLatLon: i.latlon,
	}

}
