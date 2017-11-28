package filter

import (
	"clickyab.com/crane/crane/entity"
)

// CheckAppCarrier return boolean
func CheckAppCarrier(c entity.Context, in entity.Advertise) bool {
	if len(in.Campaign().AppCarriers()) == 0 {
		return true
	}
	return hasString(in.Campaign().AppCarriers(), c.Data().PhoneData.Carrier)
}
