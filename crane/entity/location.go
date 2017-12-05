package entity

// Country is the country object
type Country struct {
	Valid bool   `json:"valid"`
	Name  string `json:"name"`
	ISO   string `json:"iso"`
}

// Province of the request
type Province struct {
	Valid bool   `json:"valid"`
	Name  string `json:"name"`
	ID    int64  `json:"id"`
}

// LatLon is the latitude longitude
type LatLon struct {
	Valid bool    `json:"valid"`
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
}

type ISP struct {
	Name  string `json:"name"`
	Valid bool   `json:"valid"`
	ID    int64  `json:"id"`
}

// Location is the location provider
type Location interface {
	// Country get the country if available
	Country() Country
	// Province get the province of request if available
	Province() Province
	// LatLon return the latitude longitude if any
	LatLon() LatLon
	// ISP return isp of the selected req
	ISP() ISP
}
