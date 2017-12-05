package filter

import (
	"clickyab.com/crane/crane/entity"
)

// CheckOS is the filter function that check for os in system
func CheckOS(c entity.Context, in entity.Advertise) bool {
	if len(in.Campaign().AllowedOS()) == 0 {
		return true
	}
	return c.OS().Valid && hasString(in.Campaign().AllowedOS(), c.OS().Name)
}
