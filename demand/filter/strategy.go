package filter

import (
	"clickyab.com/crane/demand/entity"
)

// Strategy checker
type Strategy struct {
}

// Check check if creative can be use for this impression
func (*Strategy) Check(impression entity.Context, ad entity.Creative) bool {
	return ad.Campaign().Strategy().IsSubsetOf(impression.Strategy())
}
