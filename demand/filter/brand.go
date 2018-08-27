package filter

import (
	"errors"

	"clickyab.com/crane/demand/entity"
)

// AppBrand app carrier
type AppBrand struct {
}

// Check test if campaign accept brand
func (*AppBrand) Check(c entity.Context, in entity.Creative) error {
	if hasString(true, in.Campaign().AppBrands(), c.Brand()) {
		return nil
	}
	return errors.New("APP_BRAND")

}
