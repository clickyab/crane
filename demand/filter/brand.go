package filter

import (
	"clickyab.com/crane/demand/entity"
)

// AppBrand app carrier
type AppBrand struct {
}

func (*AppBrand) Check(c entity.Context, in entity.Advertise) bool {
	return hasString(true, in.Campaign().AppBrands(), c.Brand())
}
