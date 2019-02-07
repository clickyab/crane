package campaign

import (
	"errors"

	"fmt"

	"clickyab.com/crane/demand/entity"
)

// OS Checker
type OS struct {
}

// Check is the filter function that check for os in system
func (*OS) Check(c entity.Context, in entity.Campaign) error {
	if len(in.AllowedOS()) == 0 {
		return nil
	}
	if c.OS().Valid && hasString(true, in.AllowedOS(), fmt.Sprint(c.OS().ID)) {
		return nil
	}

	return errors.New("OS")

}
