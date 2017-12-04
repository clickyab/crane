package filter

import (
	"clickyab.com/crane/crane/entity"
)

// IsNotWebMobile filter for webmobile
func IsNotWebMobile(c entity.Context, in entity.Advertise) bool {
	if c.Common().Mobile {
		return true
	}
	return !in.Campaign().WebMobile()
}

// IsWebMobile return if the campaign is ok for web mobile
func IsWebMobile(c entity.Context, in entity.Advertise) bool {
	if c.Common().Mobile {
		return in.Campaign().WebMobile()
	}

	return true
}
