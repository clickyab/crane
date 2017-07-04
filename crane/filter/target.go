package filter

import "clickyab.com/crane/crane/entity"

// Target if the campaign target and imp target is ok
func Target(imp entity.Impression, ad entity.Advertise) bool {
	for _, p := range imp.Publisher().AcceptedTargets() {
		for _, i := range ad.Campaign().Target() {
			if i == p {
				return true
			}
		}
	}
	return false
}
