package filter

import (
	"clickyab.com/crane/crane/entity"
)

// Desktop checker
type Desktop struct {
}

// Check filter network for desktop
func (*Desktop) Check(c entity.Context, in entity.Advertise) bool {
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
