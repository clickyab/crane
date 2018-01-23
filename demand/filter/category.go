package filter

import (
	"clickyab.com/crane/demand/entity"
)

// Category checker
type Category struct {
}

// Check check for category match
func (*Category) Check(context entity.Context, advertise entity.Creative) bool {
	return hasCategory(true, context.Category(), advertise.Campaign().Category())
}
