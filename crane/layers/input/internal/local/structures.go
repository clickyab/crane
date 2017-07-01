package local

import "clickyab.com/crane/crane/entity"

// Publisher Publisher
type Publisher struct {
	// Name of publisher
	FName string `json:"name"`
	// FloorCPM is the floor cpm for publisher
	FFloorCPM int64 `json:"floor_cpm"`
	// SoftFloorCPM is the soft version of floor cpm. if the publisher ahs it, then the system
	// try to use this as floor, but if this is not available, the FloorCPM is used
	FSoftFloorCPM int64 `json:"soft_floor_cpm"`
	// Attributes is the generic attribute system
	FAttributes map[string]interface{} `json:"attributes"`
	// Supplier the supplier
	FSupplier string `json:"supplier"`
	// UnderFloor asd
	FUnderFloor *bool
}

// Location Location
type Location struct {
	FCountry  entity.Country  `json:"country"`
	FProvince entity.Province `json:"province"`
	FLatLon   entity.LatLon   `json:"lat_lon"`
}

// Slot Slot
type Slot struct {
	FWidth   int    `json:"width"`
	FHeight  int    `json:"height"`
	FTrackID string `json:"track_id"`
	slotCTR  float64

	attribute map[string]interface{}
	winnerAd  interface{}
	showURL   string
}
