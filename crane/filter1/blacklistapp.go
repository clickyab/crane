package filter

import (
	"clickyab.com/crane/crane/entity"
)

// CheckAppBlackList filter blacklist
func CheckAppBlackList(c entity.Context, advertise entity.Advertise) bool {

	return !in.CampaignAppFilter.Has(false, c.GetData().App.ID)
}
