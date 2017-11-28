package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
	"clickyab.com/gad/utils"
)

// CheckAppAreaInGlob filter area in glob
func CheckAppAreaInGlob(c *builder.context, in models.AdData) bool {
	if c.GetData().CellLocation == nil || c.GetData().CellLocation.Location == "" {
		// there is no location detected
		// The campaign is regional, so ignore it
		if in.CampaignLatMap.Valid && in.CampaignLongMap.Valid && in.CampaignRadius.Valid {
			return false
		}
		return true
	}
	// The campaign is not regional, so return ok and add them to list
	if !in.CampaignLatMap.Valid || !in.CampaignLongMap.Valid || !in.CampaignRadius.Valid {
		return true
	}
	// Campaign is regional and phone is detected
	return utils.AreaInGlob(in.CampaignLatMap.Float64, in.CampaignLongMap.Float64, c.GetData().CellLocation.Lat, c.GetData().CellLocation.Lon, in.CampaignRadius.Float64)
}
