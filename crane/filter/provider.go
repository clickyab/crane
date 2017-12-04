package filter

import (
	"clickyab.com/crane/crane/entity"
)

//CheckProvder find provider client in campaign
func CheckProvder(c entity.Context, in entity.Advertise) bool {
	if c.Data().PhoneData.Network == "" {
		return len(in.Campaign().NetProvider()) == 0
	}
	return hasString(in.Campaign().NetProvider(), c.Data().PhoneData.Network)
}
