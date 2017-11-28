package filter

import (
	"clickyab.com/crane/crane/entity"
)

// Country check if the ad country is matched or not
func Country(impression entity.Context, ad entity.Advertise) bool {
	adCountry := ad.Campaign().Country()
	if len(adCountry) == 0 {
		// the ad has no country attached. so pass it
		return true
	}
	// the ad has country
	country := impression.Location().Country()
	if !country.Valid {
		// ad has country but the imp has no country, so ignore it
		return false
	}

	return hasString(adCountry, country.Name)
}
