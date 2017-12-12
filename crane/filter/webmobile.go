package filter

import (
	"clickyab.com/crane/crane/entity"
)

// WebMobile checker
type WebMobile struct {
}

// Check return if the campaign is ok for web mobile
func (*WebMobile) Check(c entity.Context, in entity.Advertise) bool {
	if c.IsMobile() {
		return in.Campaign().WebMobile()
	}

	return true
}
