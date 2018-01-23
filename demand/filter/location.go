package filter

import (
	"math"

	"clickyab.com/crane/demand/entity"
)

// areaInGlob is a helper fuunction to handle check point in a globe
func areaInGlob(lat, lon, centerLat, centerLon, radius float64) bool {
	var ky = 40000.0 / 360.0
	var kx = math.Cos(math.Pi*centerLat/180.0) * ky
	dx := math.Abs(centerLon-lon) * kx
	dy := math.Abs(centerLat-lat) * ky
	return math.Sqrt(dx*dx+dy*dy) <= radius
}

// AreaInGlob is a mobile checker for area, if location is available
type AreaInGlob struct {
}

// Check filter area in glob
func (*AreaInGlob) Check(c entity.Context, in entity.Creative) bool {
	b, lat, lon, radius := in.Campaign().LatLon()
	ll := c.Location().LatLon()
	if !ll.Valid {
		// there is no location detected
		// if the campaign is regional, ignore it
		return !b
		// no location and no regional campaign so be it!
	}
	// The campaign is not regional, so return ok and add them to list
	if !b {
		return true
	}
	// Campaign is regional and phone is detected

	return areaInGlob(lat, lon, ll.Lat, ll.Lon, radius)
}
