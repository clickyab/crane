package filter

import (
	"clickyab.com/crane/demand/entity"
)

// OS Checker
type OS struct {
}

// Check is the filter function that check for os in system
func (*OS) Check(c entity.Context, in entity.Creative) bool {
	if len(in.Campaign().AllowedOS()) == 0 {
		return true
	}
	return c.OS().Valid && hasString(true, in.Campaign().AllowedOS(), c.OS().Name)
}
