package filter

import (
	"errors"

	"fmt"

	"clickyab.com/crane/demand/entity"
)

// Province checker
type Province struct {
}

//Check find province client in campaign
func (*Province) Check(c entity.Context, in entity.Creative) error {
	if c.Location().Province().Name == "" {
		if len(in.Campaign().Province()) == 0 {
			return nil
		}
		return errors.New("province is unknown")
	}

	// the 1 means for iran
	if c.Location().Country().ISO == "IR" {
		if hasString(true, in.Campaign().Province(), "1") {
			return nil
		}
	}

	if hasString(true, in.Campaign().Province(), fmt.Sprint(c.Location().Province().ID)) {
		return nil
	}
	return errors.New("province is not allowed")
}
