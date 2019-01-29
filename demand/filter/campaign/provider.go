package campaign

import (
	"errors"

	"clickyab.com/crane/demand/entity"
)

// ConnectionType check app network
type ConnectionType struct {
}

// Check test if campaign accept provider
func (*ConnectionType) Check(c entity.Context, in entity.Campaign) error {
	if hasInt(true, in.ConnectionType(), int(c.ConnectionType())) {
		return nil
	}
	return errors.New("CONNECTION")
}
