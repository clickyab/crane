package filter

import (
	"clickyab.com/crane/demand/entity"
)

// Desktop checker
type Desktop struct {
}

// Check filter network for desktop
func (*Desktop) Check(c entity.Context, in entity.Creative) bool {
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
