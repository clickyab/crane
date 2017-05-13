package ip2location

// IP2lData struct
type IP2lData struct {
	CountryShort string `json:"country_short"`
	CountryLong  string `json:"country_long"`
	Region       string `json:"region"`
	City         string `json:"city"`
}

// IP2Location is the local use of this function
func IP2Location(ip string) IP2lData {
	rec := GetAll(ip)
	return IP2lData{
		CountryShort: rec.Country_short,
		CountryLong:  rec.Country_long,
		Region:       rec.Region,
		City:         rec.City,
	}
}
