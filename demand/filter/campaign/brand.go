package campaign

import (
	"errors"

	"clickyab.com/crane/demand/entity"
)

// AppBrand app carrier
type AppBrand struct {
}

// Check test if campaign accept brand
func (*AppBrand) Check(c entity.Context, in entity.Campaign) error {
	if hasString(true, in.AppBrands(), c.Brand()) {
		return nil
	}
	return errors.New("APP_BRAND")

}
