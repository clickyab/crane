package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
)

// CheckAppCarrier return boolean
func CheckAppCarrier(c *builder.context, in models.AdData) bool {
	return in.Campaign.CampaignAppsCarriers.Has(true, c.GetData().PhoneData.CarrierID)
}
