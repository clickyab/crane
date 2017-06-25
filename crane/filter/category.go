package filter

import (
	"clickyab.com/crane/crane/entity"
)

// Category check for category match
func Category(impression entity.Impression, advertise entity.Advertise) bool {
	whitelist := advertise.Campaign().Category()
	elems := impression.Category()

	return hasCategory(elems, whitelist)
}
