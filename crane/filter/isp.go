package filter

import (
	"fmt"

	"clickyab.com/crane/crane/entity"
)

//CheckISP find isp
func CheckISP(c entity.Context, in entity.Advertise) bool {
	isp := c.Location().ISP().ID
	if isp == 0 {
		return len(in.Campaign().ISP()) == 0
	}
	return hasString(true, in.Campaign().ISP(), fmt.Sprint(isp))
}
