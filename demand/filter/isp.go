package filter

import (
	"fmt"

	"clickyab.com/crane/demand/entity"
)

// ISP Checker
type ISP struct {
}

//Check find isp
func (*ISP) Check(c entity.Context, in entity.Creative) bool {
	isp := c.Location().ISP().ID
	if isp == 0 {
		return len(in.Campaign().ISP()) == 0
	}
	return hasString(true, in.Campaign().ISP(), fmt.Sprint(isp))
}
