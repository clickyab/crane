package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
)

// CheckOS is the filter function that check for os in system
func CheckOS(c *builder.context, in models.AdData) bool {
	return in.CampaignPlatforms.Has(true, c.GetCommon().PlatformID)
}
