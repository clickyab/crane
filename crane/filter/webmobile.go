package filter

import (
	"clickyab.com/crane/crane/entity"
)

// IsWebMobile return if the campaign is ok for web mobile
func IsWebMobile(c entity.Context, in entity.Advertise) bool {
	if c.IsMobile() {
		return in.Campaign().WebMobile()
	}

	return true
}
