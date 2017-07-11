package native

import (
	"net"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/layers/input/internal/local"
	"github.com/clickyab/services/ip2location"
)

func (i *imp) TrackID() string {
	return i.fTrackID
}

func (i *imp) ClientID() string {
	return i.fClientID
}

func (i *imp) IP() net.IP {
	return i.fIP
}

func (i *imp) UserAgent() string {
	return i.fUA
}

func (i *imp) Source() entity.Publisher {
	return i.fPub
}

func (i *imp) Location() entity.Location {
	return i.fLocation
}

func (i *imp) OS() entity.OS {
	return i.fOS
}

func (i *imp) Slots() []entity.Slot {
	if i.nDum == nil {
		i.nDum = make([]entity.Slot, len(i.fSlots))
		for j := range i.fSlots {
			i.nDum[j] = i.fSlots[j]
		}
	}
	return i.nDum
}

func (i *imp) Category() []entity.Category {
	return i.fCategories
}

func (i *imp) Publisher() entity.Publisher {
	return i.fPub
}

func (i *imp) Protocol() string {
	return i.fprotocol
}

func (i *imp) extractData() {
	d := ip2location.IP2Location(i.fIP.String())
	i.fLocation = &local.Location{
		FCountry: entity.Country{
			Name:  d.CountryLong,
			ISO:   d.CountryShort,
			Valid: d.CountryLong != "-",
		},

		FProvince: entity.Province{
			Valid: d.Region != "-",
			Name:  d.Region,
		},

		FLatLon: i.latlon,
	}

}
