package filter

import (
	"clickyab.com/crane/demand/entity"
)

// AppProvider check app network
type AppProvider struct {
}

func (*AppProvider) Check(c entity.Context, in entity.Advertise) bool {
	provider, _ := c.Network()
	return hasString(true, in.Campaign().NetProvider(), provider)
}
