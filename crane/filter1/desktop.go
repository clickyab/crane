package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
)

// CheckDesktopNetwork filter network for desktop
func CheckDesktopNetwork(c *builder.context, in models.AdData) bool {
	if in.CampaignWeb == 1 {
		if in.CampaignWebMobile == 0 {
			return !c.GetCommon().Mobile
		}
	} else if in.CampaignWeb == 0 {
		if in.CampaignWebMobile == 1 {
			return c.GetCommon().Mobile
		}
	}
	return true
}
