package filter

import (
	"errors"
	"fmt"

	"clickyab.com/crane/demand/entity"
)

// ISP Checker
type ISP struct {
}

//Check find isp
func (*ISP) Check(c entity.Context, in entity.Creative) error {
	isp := c.Location().ISP().ID
	if isp == 0 {
		if len(in.Campaign().ISP()) == 0 {
			return nil
		}
		return errors.New("isp filter not met")
	}
	if hasString(true, in.Campaign().ISP(), fmt.Sprint(isp)) {
		return nil
	}
	return errors.New("isp filter not met")
}
