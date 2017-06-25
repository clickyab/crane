package filter

import (
	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/reducer"
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
			if t.IsSizeAllowed(advertise.Width(), advertise.Height()) && i.IsSizeAllowed(advertise.Width(), advertise.Height()) {
				return true
			}
		}

		return false
	}

}
