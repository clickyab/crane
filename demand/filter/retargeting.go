package filter

import (
	"fmt"

	"clickyab.com/crane/demand/entity"
)

// ReTargeting checker
type ReTargeting struct {
}

// Check check if creative can be use for this impression
func (*ReTargeting) Check(impression entity.Context, ad entity.Creative) error {
	if len(ad.Campaign().ReTargeting()) == 0 {
		return nil
	}

	for _, v := range ad.Campaign().ReTargeting() {
		if _, ok := impression.User().List()[v]; ok {
			return nil
		}
	}

	return fmt.Errorf("retargeting failed")
}
