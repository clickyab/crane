package filter

import (
	"clickyab.com/crane/demand/entity"
)

// AppCarrier app carrier
type AppCarrier struct {
}

func (*AppCarrier) Check(c entity.Context, in entity.Advertise) bool {
	carrString, _ := c.Carrier()
	return hasString(true, in.Campaign().AppCarriers(), carrString)
}
