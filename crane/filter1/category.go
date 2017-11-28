package filter

import (
	"clickyab.com/crane/crane/entity"
)

// Category check for category match
func Category(context entity.Context, advertise entity.Advertise) bool {
	whitelist := advertise.Campaign().Category()
	elems := context.Category()

	return hasCategory(elems, whitelist)
}
