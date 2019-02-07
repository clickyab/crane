package campaign

import (
	"errors"

	"clickyab.com/crane/demand/entity"
)

// Category filter for requests
type Category struct {
}

// Check iab category
func (*Category) Check(c entity.Context, in entity.Campaign) error {
	if len(in.Category()) == 0 {
		return nil
	}
	f := make(map[string]bool)
	for _, v := range c.Category() {
		f[string(v)] = true
	}

	for _, v := range in.Category() {
		if f[string(v)] {

			return nil
		}
	}

	return errors.New("IAB")

}
