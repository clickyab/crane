package filter

import (
	"clickyab.com/crane/crane/entity"
)

//CheckISP find isp
func CheckISP(c entity.Context, in entity.Advertise) bool {
	if c.Isp() == "" {
		return len(in.Campaign().Isp()) == 0
	}
	return hasString(in.Campaign().Isp(), c.Isp())
}
