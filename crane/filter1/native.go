package filter

import (
	"clickyab.com/crane/crane/entity"
)

// IsNativeAd tells if an ad is native
func IsNativeAd(c entity.Context, in entity.Advertise) bool {
	return in.Type() == entity.AdTypeNative
}
