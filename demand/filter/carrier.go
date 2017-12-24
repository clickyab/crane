package filter

import (
	"clickyab.com/crane/demand/entity"
)

// AppCarrier app carrier
type AppCarrier struct {
}

// Check test if campaign accept carrier
func (*AppCarrier) Check(c entity.Context, in entity.Advertise) bool {
	return hasString(true, in.Campaign().AppCarriers(), c.Carrier())
}
