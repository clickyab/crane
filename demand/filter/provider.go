package filter

import (
	"clickyab.com/crane/demand/entity"
)

// AppProvider check app network
type AppProvider struct {
}

// Check test if campaign accept provider
func (*AppProvider) Check(c entity.Context, in entity.Creative) bool {
	return hasString(true, in.Campaign().NetProvider(), c.Network())
}
