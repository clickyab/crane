package filter

import (
	"errors"

	"clickyab.com/crane/demand/entity"
)

// ConnectionType check app network
type ConnectionType struct {
}

// Check test if campaign accept provider
func (*ConnectionType) Check(c entity.Context, in entity.Creative) error {
	if hasInt(true, in.Campaign().ConnectionType(), int(c.ConnectionType())) {
		return nil
	}
	return errors.New("CONNECTION")
}
