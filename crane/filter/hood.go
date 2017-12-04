package filter

import (
	"clickyab.com/crane/crane/entity"
)

// CheckAppHood return boolean
func CheckAppHood(c entity.Context, in entity.Advertise) bool {
	if c.Data().CellLocation == nil {
		return len(in.Campaign().Hoods()) == 0
	}
	return hasString(in.Campaign().Hoods(), c.Data().CellLocation.Location)
}
