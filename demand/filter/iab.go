package filter

import (
	"fmt"
	"strings"

	"clickyab.com/crane/demand/entity"
)

// Category filter for requests
type Category struct {
}

// Check iab category
func (*Category) Check(c entity.Context, in entity.Creative) error {
	if len(in.Campaign().Category()) == 0 {
		return nil
	}
	f := make(map[string]bool)
	for _, v := range c.Category() {
		f[strings.Split(string(v), "-")[0]] = true
	}
	for _, v := range in.Campaign().Category() {
		if f[string(v)] {
			return nil
		}
	}
	return fmt.Errorf("category filter not met")
}
