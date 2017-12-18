package builder

import "math"

// AreaInGlob is a helper function to handle check point in a globe
func AreaInGlob(lat, lon, centerLat, centerLon, radius float64) bool {
	var ky = 40000.0 / 360.0
	var kx = math.Cos(math.Pi*centerLat/180.0) * ky
	dx := math.Abs(centerLon-lon) * kx
	dy := math.Abs(centerLat-lat) * ky
	return math.Sqrt(dx*dx+dy*dy) <= radius
}
