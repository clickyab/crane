package entity

// Country is the country object
type Country struct {
	Valid bool
	Name  string
	ISO   string
}

// Province of the request
type Province struct {
	Valid bool
	Name  string
}

// LatLon is the latitude longitude
type LatLon struct {
	Valid    bool
	Lat, Lon float64
}

// Location is the location provider
type Location interface {
	// Country get the country if available
	Country() Country
	// Province get the province of request if available
	Province() Province
	// LatLon return the latitude longitude if any
	LatLon() LatLon
}
