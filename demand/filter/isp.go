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
	if hasString(true, in.Campaign().ISP(), fmt.Sprint(c.Location().ISP().ID)) {
		return nil
	}
	return errors.New("ISP")
}
