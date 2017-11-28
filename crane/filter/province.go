package filter

import (
	"clickyab.com/crane/crane/entity"
)

// Province check if the ad accept this province or not
func Province(impression entity.Context, ad entity.Advertise) bool {
	elem := ad.Campaign().Province()
	if len(elem) == 0 {
		// ad has no province attach to it
		return true
	}
	province := impression.Location().Province()
	if !province.Valid {
		// ad need province but we can not detect it on imp not pass it
		return false
	}

	return hasString(elem, province.Name)
}
