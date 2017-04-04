package location

import (
	"entity"
	"net"
	"services/gmaps"
	"services/location/internal/models"

	"github.com/Sirupsen/logrus"
)

type data struct {
	countryCode        string
	iploc              *models.IP2Location
	mcc, mnc, lac, cid int64
	latlon             *entity.LatLon
	country            *entity.Country
	province           *entity.Province
}

// Country get the country if available
func (d *data) Country() entity.Country {
	if d.country != nil {
		return *d.country
	}
	d.country = &entity.Country{}
	m := models.NewManager()
	cCode := d.countryCode
	if cCode == "" && d.iploc != nil {
		cCode = d.iploc.CountryCode.String
	}

	if cCode != "" {
		c, err := m.GetCountry(cCode)
		if err == nil {
			d.country.Valid = true
			d.country.ID = c.ID
			d.country.ISO = c.Iso
			d.country.Name = c.Name
		}
	}

	return *d.country
}

// Province get the province of request if available
func (d *data) Province() entity.Province {
	if d.province != nil {
		return *d.province
	}
	d.province = &entity.Province{}
	m := models.NewManager()
	pCode := ""
	if d.iploc != nil {
		pCode = d.iploc.RegionName.String
	}

	if pCode != "" {
		c, err := m.GetProvince(pCode)
		if err == nil {
			d.province.Valid = true
			d.province.ID = c.ID
			d.province.Name = c.Name
		}
	}

	return *d.province
}

// LatLon return the latitude longitude if any
func (d *data) LatLon() entity.LatLon {
	if d.latlon == nil {
		if d.mcc == 0 && d.mnc == 0 && d.lac == 0 && d.cid == 0 {
			d.latlon = &entity.LatLon{}
		} else {
			lat, lon, err := gmaps.LockUp(d.mcc, d.mnc, d.lac, d.cid)
			if err != nil {
				d.latlon = &entity.LatLon{}
			} else {
				d.latlon = &entity.LatLon{
					Valid: true,
					Lat:   lat,
					Lon:   lon,
				}
			}
		}
	}
	return *d.latlon
}

// Provider is the service entry for ip to system location
func Provider(ip net.IP, countryCode string, mcc, mnc, lac, cid int64) entity.Location {
	m := models.NewManager()
	loc, err := m.GetLocation(ip)
	if err != nil {
		logrus.Debug(err)
	}
	return &data{
		countryCode: countryCode,
		iploc:       loc,
		mcc:         mcc,
		mnc:         mnc,
		lac:         lac,
		cid:         cid,
	}
}

// ProviderSimple is the shorter version of the ip to location
func ProviderSimple(ip net.IP, countrCode string) entity.Location {
	return Provider(ip, countrCode, 0, 0, 0, 0)
}
