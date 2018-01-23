package filter

import (
	"clickyab.com/crane/demand/entity"
)

// Desktop checker
type Desktop struct {
}

// Check filter network for desktop
func (*Desktop) Check(c entity.Context, in entity.Creative) bool {
	if c.IsMobile() {
		return in.Campaign().WebMobile()
	}
	return in.Campaign().Web()
}
