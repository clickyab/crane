package filter

import (
	"clickyab.com/crane/crane/entity"
)

// WebNetwork is the web network checker
type WebNetwork struct {
}

// Check filter network for campaigns
func (*WebNetwork) Check(c entity.Context, in entity.Advertise) (b bool) {
	if in.Campaign().Target() == entity.TargetWeb {
		return in.Campaign().Web() || in.Campaign().WebMobile()
	}
	return in.Campaign().Target() == entity.TargetWeb
}

// AppNetwork checker
type AppNetwork struct {
}

// Check filter network for campaigns
func (*AppNetwork) Check(c entity.Context, in entity.Advertise) bool {
	return in.Campaign().Target() == entity.TargetApp
}

// NativeNetwork checker
type NativeNetwork struct {
}

// Check filter network for native
func (*NativeNetwork) Check(c entity.Context, in entity.Advertise) bool {
	return in.Campaign().Target() == entity.TargetNative
}

// VastNetwork checker
type VastNetwork struct {
}

// Check filter network for native
func (*VastNetwork) Check(c entity.Context, in entity.Advertise) bool {
	return in.Campaign().Target() == entity.TargetVast
}
