package filter

import (
	"clickyab.com/crane/demand/entity"
)

var videoValidSizes = []int{3, 4, 9, 16, 14}

// WebSize checker
type WebSize struct {
}

// Check check if the banner size exists in the request
func (*WebSize) Check(c entity.Context, in entity.Advertise) bool {
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

// AppSize checker
type AppSize struct {
}

// Check check if the banner size exists in the request
func (*AppSize) Check(c entity.Context, in entity.Advertise) bool {
	// TODO : fix the dynamic click and remove the second condition
	if in.Type() == entity.AdTypeVideo || in.Type() == entity.AdTypeDynamic {
		return false
	}
	for _, seat := range c.Seats() {
		if seat.Size() == in.Size() {
			return true
		}
	}
	return false
}
