package local

import "clickyab.com/crane/crane/entity"

// Location Location
type Location struct {
	FCountry  entity.Country  `json:"country"`
	FProvince entity.Province `json:"province"`
	FLatLon   entity.LatLon   `json:"lat_lon"`
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
