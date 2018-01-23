package filter

import (
	"clickyab.com/crane/demand/entity"
)

// AppBrand app carrier
type AppBrand struct {
}

// Check test if campaign accept brand
func (*AppBrand) Check(c entity.Context, in entity.Creative) bool {
	return hasString(true, in.Campaign().AppBrands(), c.Brand())
}
