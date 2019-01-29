package creative

import (
	"errors"

	"fmt"

	"clickyab.com/crane/demand/entity"
)

// OS Checker
type OS struct {
}

// Check is the filter function that check for os in system
func (*OS) Check(c entity.Context, in entity.Creative) error {
	if len(in.Campaign().AllowedOS()) == 0 {
		return nil
	}
	if c.OS().Valid && hasString(true, in.Campaign().AllowedOS(), fmt.Sprint(c.OS().ID)) {
		return nil
	}

	return errors.New("OS")

}
