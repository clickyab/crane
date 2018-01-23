package filter

import (
	"clickyab.com/crane/demand/entity"
)

// WebNetwork is the web network checker
type WebNetwork struct {
}

// Check filter network for campaigns
func (*WebNetwork) Check(c entity.Context, in entity.Creative) (b bool) {
	if in.Campaign().Target() == entity.TargetWeb {
		return in.Campaign().Web() || in.Campaign().WebMobile()
	}
	return in.Campaign().Target() == entity.TargetWeb
}

// AppNetwork checker
type AppNetwork struct {
}

// Check filter network for campaigns
func (*AppNetwork) Check(c entity.Context, in entity.Creative) bool {
	return in.Campaign().Target() == entity.TargetApp
}

// NativeNetwork checker
type NativeNetwork struct {
}

// Check filter network for native
func (*NativeNetwork) Check(c entity.Context, in entity.Creative) bool {
	return in.Campaign().Target() == entity.TargetNative
}

// VastNetwork checker
type VastNetwork struct {
}

// Check filter network for native
func (*VastNetwork) Check(c entity.Context, in entity.Creative) bool {
	// TODO : the following line is correct. but, since we use an invalid form of ads in our system, we should comment it
	//return in.Campaign().Target() == entity.TargetVast
	if in.Campaign().Target() != entity.TargetVast {
		// TODO : remove it when the new console is awaken!
		if in.Size() != 9 && in.Campaign().Target() != entity.TargetWeb { // there is a fucking decision to show web size 9 in vast network.
			return false
		}
	}
	return true
}
