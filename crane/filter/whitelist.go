package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
)

// CheckWhiteList return boolean
func CheckWhiteList(c *builder.context, in models.AdData) bool {
	return in.CampaignPlacement.Has(true, c.GetData().Website.WID)
}

// CheckAppWhiteList return boolean
func CheckAppWhiteList(c *builder.context, in models.AdData) bool {
	return in.CampaignApp.Has(true, c.GetData().App.ID)
}
