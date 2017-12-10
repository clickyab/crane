package filter

import (
	"clickyab.com/crane/crane/entity"
)

var videoValidSizes = []int{3, 4, 9, 16, 14}

// CheckWebSize check if the banner size exists in the request
func CheckWebSize(c entity.Context, in entity.Advertise) bool {
	if in.Type() == entity.AdTypeVideo {
		for _, seat := range c.Seats() {
			if hasInt(false, videoValidSizes, seat.Size()) {
				return true
			}
		}
		return false
	}

	for _, seat := range c.Seats() {
		if seat.Size() == in.Size() {
			return true
		}
	}
	return false
}
