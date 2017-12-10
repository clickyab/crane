package filter

import (
	"clickyab.com/crane/crane/entity"
)

// Category check for category match
func Category(context entity.Context, advertise entity.Advertise) bool {
	return hasCategory(true, context.Category(), advertise.Campaign().Category())
}
