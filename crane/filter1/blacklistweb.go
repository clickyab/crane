package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
)

// CheckWebBlackList filter blacklist
func CheckWebBlackList(c *builder.context, in models.AdData) bool {
	return !in.CampaignWebsiteFilter.Has(false, c.GetData().Website.WID)
}
