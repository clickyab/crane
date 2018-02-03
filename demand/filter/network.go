package filter

import "clickyab.com/crane/demand/entity"

// Network checker
type Network struct {
}

// Check check if creative can be use for this network
func (*Network) Check(impression entity.Context, ad entity.Creative) bool {
	return impression.Publisher().Supplier().Strategy().Has(ad.Campaign().Strategy())
}
