package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
)

// CheckVastNetwork filter vast video so vast (campaign)
func CheckVastNetwork(c *builder.context, in models.AdData) bool {
	if in.CampaignNetwork != 2 && in.AdType == 3 {
		return false
	}
	return true
}

// CheckVastOtherNetwork filter vast video so to not be shown in other select ads like (web,native,app,...)
func CheckVastOtherNetwork(c *builder.context, in models.AdData) bool {
	if in.CampaignNetwork == 2 {
		return false
	}
	return true
}
