package filter

import (
	"clickyab.com/crane/demand/entity"
)

// AppBrand app carrier
type AppBrand struct {
}

func (*AppBrand) Check(c entity.Context, in entity.Advertise) bool {
	brand, _ := c.Brand()
	return hasString(true, in.Campaign().AppBrands(), brand)
}
