package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
)

// CheckAppHood return boolean
func CheckAppHood(c *builder.context, in models.AdData) bool {
	if c.GetData().CellLocation == nil {
		return in.Campaign.CampaignHoods == ""
	}
	return in.Campaign.CampaignHoods.Has(true, c.GetData().CellLocation.NeighborhoodsID)
}
