package local

import "clickyab.com/crane/crane/entity"

// Location Location
type Location struct {
	TheCountry  entity.Country  `json:"country"`
	TheProvince entity.Province `json:"province"`
	TheLatLon   entity.LatLon   `json:"latlon"`
}

// Publisher Publisher
type Publisher struct {
}

// Slot Slot
type Slot struct {
}
