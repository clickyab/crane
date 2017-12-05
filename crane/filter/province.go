package filter

import (
	"clickyab.com/crane/crane/entity"
)

//CheckProvince find province client in campaign
func CheckProvince(c entity.Context, in entity.Advertise) bool {
	if c.Location().Province().Name == "" {
		return len(in.Campaign().Province()) == 0
	}
	// The 1 means iran. watch for it please!
	return hasString(in.Campaign().Province(), c.Location().Province().Name)
}
