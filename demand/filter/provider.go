package filter

import (
	"clickyab.com/crane/demand/entity"
)

// ConnectionType check app network
type ConnectionType struct {
}

// Check test if campaign accept provider
func (*ConnectionType) Check(c entity.Context, in entity.Creative) bool {
	return hasInt(true, in.Campaign().ConnectionType(), c.ConnectionType())
}
