package filter

import (
	"clickyab.com/crane/demand/entity"
)

// WebMobile checker
type WebMobile struct {
}

// Check return if the campaign is ok for web mobile
func (*WebMobile) Check(c entity.Context, in entity.Creative) bool {
	if c.IsMobile() {
		return in.Campaign().WebMobile()
	}

	return true
}
