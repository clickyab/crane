package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/crane/crane/entity"
)

// CheckAppAreaInGlob filter area in glob
func CheckAppAreaInGlob(c entity.Context, in entity.Advertise) bool {
	lat, lon, rad := in.Campaign().LanLon()
	if c.Data().CellLocation == nil || c.Data().CellLocation.Location == "" {
		// there is no location detected
		// The campaign is regional, so ignore it

		if lat != 0 && lon != 0 && rad != 0 {
			return false
		}
		return true
	}
	// The campaign is not regional, so return ok and add them to list
	if lat == 0 || lon == 0 || rad == 0 {
		return true
	}
	// Campaign is regional and phone is detected
	return builder.AreaInGlob(lat, lon, c.Data().CellLocation.Lat, c.Data().CellLocation.Lon, rad)
}
