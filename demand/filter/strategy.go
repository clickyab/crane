package filter

import (
	"fmt"

	"clickyab.com/crane/demand/entity"
)

// Strategy checker
type Strategy struct {
}

// Check check if creative can be use for this impression
func (*Strategy) Check(impression entity.Context, ad entity.Creative) error {
	if ad.Campaign().Strategy().IsSubsetOf(impression.Strategy()) {
		return nil
	}
	return fmt.Errorf("supplier strategy is %d but campaign want %d ",
		impression.Strategy(), ad.Campaign().Strategy())
}
