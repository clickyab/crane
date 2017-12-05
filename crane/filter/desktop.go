package filter

import (
	"clickyab.com/crane/crane/entity"
)

// CheckDesktopNetwork filter network for desktop
func CheckDesktopNetwork(c entity.Context, in entity.Advertise) bool {
	if in.Campaign().Web() {
		if !in.Campaign().WebMobile() {
			return !c.IsMobile()
		}
	} else {
		if in.Campaign().WebMobile() {
			return c.IsMobile()
		}
	}
	return true
}
