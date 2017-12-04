package filter

import (
	"clickyab.com/crane/crane/entity"
)

// IsWebNetwork filter network for campaigns
func IsWebNetwork(c entity.Context, in entity.Advertise) bool {
	if in.Campaign().Target() == entity.TargetWeb {
		return in.Campaign().Web() || in.Campaign().WebMobile()
	}
	return in.Campaign().Target() == entity.TargetWeb || in.Campaign().Target() == entity.TargetVast
}

// IsAppNetwork filter network for campaigns
func IsAppNetwork(c entity.Context, in entity.Advertise) bool {
	return in.Campaign().Target() == entity.TargetApp
}

// IsNativeNetwork filter network for native
func IsNativeNetwork(c entity.Context, in entity.Advertise) bool {
	return in.Campaign().Target() == entity.TargetNative
}
