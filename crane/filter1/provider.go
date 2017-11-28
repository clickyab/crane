package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
)

//CheckProvder find provider client in campaign
func CheckProvder(c *builder.context, in models.AdData) bool {
	if c.GetData().PhoneData.NetworkID == 0 {
		return len(in.CampaignNetProvider) == 0
	}
	return in.CampaignNetProvider.Has(true, c.GetData().PhoneData.NetworkID)
}
