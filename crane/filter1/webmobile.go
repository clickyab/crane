package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
)

// IsNotWebMobile filter for webmobile
func IsNotWebMobile(c *builder.context, in models.AdData) bool {
	if c.GetCommon().Mobile {
		return true
	}
	return in.CampaignWebMobile == 0
}

// IsWebMobile return if the campaign is ok for web mobile
func IsWebMobile(c *builder.context, in models.AdData) bool {
	if c.GetCommon().Mobile {
		return in.CampaignWebMobile == 1
	}

	return true
}
