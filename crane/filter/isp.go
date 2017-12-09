package filter

import (
	"clickyab.com/crane/crane/entity"
)

//CheckISP find isp
func CheckISP(c entity.Context, in entity.Advertise) bool {
	if c.ISP() == "" {
		return len(in.Campaign().ISP()) == 0
	}
	return hasString(in.Campaign().ISP(), c.ISP())
}
