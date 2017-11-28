package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
)

// CheckWebCategory is the filter for category
func CheckWebCategory(c *builder.context, in models.AdData) bool {
	return in.CampaignCat.Match(true, c.GetData().Website.WCategories)
}

// CheckAppCategory is the filter for category
func CheckAppCategory(c *builder.context, in models.AdData) bool {
	return in.CampaignCat.Match(true, c.GetData().App.Appcat)
}
