package filter

import (
	"clickyab.com/exchange/crane/entity"
	"clickyab.com/exchange/crane/reducer"
)

var (
	// AppSize is a function to handle the app size
	AppSize = createSizeFilter(entity.TargetApp)
	// WebSize is the function to handle web size
	WebSize = createSizeFilter(entity.TargetWeb)
	// VastSize is the function to handle vast size
	VastSize = createSizeFilter(entity.TargetVast)
)

func createSizeFilter(t entity.Target) reducer.FilterFunc {
	return func(impression entity.Impression, advertise entity.Advertise) bool {
		if impType := impression.Source().AcceptedTarget(); impType != t {
			// if the impression is not this type, then pass it by
			return true
		}
		for _, i := range impression.Slots() {
			if hasInt(i.AllowedSize(), advertise.Size()) {
				return true
			}
		}

		return false
	}

}
