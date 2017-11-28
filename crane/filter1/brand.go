package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
)

// CheckAppBrand return boolean
func CheckAppBrand(c *builder.context, in models.AdData) bool {
	return in.Campaign.CampaignAppBrand.Has(true, c.GetData().PhoneData.BrandID)
}
