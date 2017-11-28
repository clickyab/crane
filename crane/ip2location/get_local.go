package ip2location

import "github.com/clickyab/services/ip2location"

// IP2lData struct
type IP2lData struct {
	CountryShort string `json:"country_short"`
	CountryLong  string `json:"country_long"`
	Region       string `json:"region"`
	City         string `json:"city"`
	ISP          string `json:"isp"`
}

// IP2Location is the local use of this function
func IP2Location(ip string) IP2lData {
	rec := ip2location.GetAll(ip)
	return IP2lData{
		CountryShort: rec.CountryShort,
		CountryLong:  rec.CountryLong,
		Region:       rec.Region,
		City:         rec.City,
		ISP:          rec.Isp,
	}
}
