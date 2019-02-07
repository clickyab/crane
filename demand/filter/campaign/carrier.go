package campaign

import (
	"errors"

	"clickyab.com/crane/demand/entity"
)

// AppCarrier app carrier
type AppCarrier struct {
}

// Check test if campaign accept carrier
func (*AppCarrier) Check(c entity.Context, in entity.Campaign) error {
	if hasString(true, in.AppCarriers(), c.Carrier()) {
		return nil
	}
	return errors.New("APP_CARRIER")

}
