package filter

import (
	"clickyab.com/crane/crane/entity"
)

// CheckAppBrand return boolean
func CheckAppBrand(c entity.Context, in entity.Advertise) bool {
	if len(in.Campaign().AppBrands()) == 0 {
		return true
	}
	return hasInt64(in.Campaign().AppBrands(), c.Data().PhoneData.BrandID)
}
