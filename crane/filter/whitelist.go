package filter

import (
	"clickyab.com/crane/crane/entity"
)

// CheckWhiteList return boolean
func CheckWhiteList(c entity.Context, in entity.Advertise) bool {
	return hasString(in.Campaign().WhiteListPublisher(), c.Publisher().Name())
}
