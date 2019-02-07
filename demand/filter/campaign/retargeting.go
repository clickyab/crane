package campaign

import (
	"fmt"

	"clickyab.com/crane/demand/entity"
)

// ReTargeting checker
type ReTargeting struct {
}

// Check check if creative can be use for this impression
func (*ReTargeting) Check(impression entity.Context, ad entity.Campaign) error {
	if len(ad.ReTargeting()) == 0 {
		return nil
	}
	fmt.Println(fmt.Sprintf("FILTER %v, %v, %v", impression.User().ID(), impression.User().List(), ad.ReTargeting()))
	if len(impression.User().List()) == 0 {
		return fmt.Errorf("retargeting failed")
	}
	for _, v := range ad.ReTargeting() {
		if _, ok := impression.User().List()[v]; ok {
			return nil
		}
	}

	return fmt.Errorf("retargeting failed")
}
