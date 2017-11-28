package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
)

//CheckProvince find province client in campaign
func CheckProvince(c *builder.context, in models.AdData) bool {
	if c.GetCommon().ProvinceID == 0 {
		return len(in.CampaignGeos) == 0
	}
	// The 1 means iran. watch for it please!
	return in.CampaignGeos.Has(true, c.GetCommon().ProvinceID, 1)
}
