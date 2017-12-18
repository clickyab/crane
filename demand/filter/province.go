package filter

import (
	"clickyab.com/crane/demand/entity"
)

// Province checker
type Province struct {
}

//Check find province client in campaign
func (*Province) Check(c entity.Context, in entity.Advertise) bool {
	if c.Location().Province().Name == "" {
		return len(in.Campaign().Province()) == 0
	}
	// The 1 means iran. watch for it please!
	return hasString(true, in.Campaign().Province(), c.Location().Province().Name)
}
