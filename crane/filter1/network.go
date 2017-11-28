package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
)

// IsWebNetwork filter network for campaigns
func IsWebNetwork(c *builder.context, in models.AdData) bool {
	if in.CampaignNetwork == 0 {
		return in.CampaignWeb == 1 || in.CampaignWebMobile == 1
	}
	return in.CampaignNetwork == 0 || in.CampaignNetwork == 2
}

// IsAppNetwork filter network for campaigns
func IsAppNetwork(c *builder.context, in models.AdData) bool {
	return in.CampaignNetwork == 1
}

// IsNativeNetwork filter network for native
func IsNativeNetwork(c *builder.context, in models.AdData) bool {
	return in.CampaignNetwork == 3
}
