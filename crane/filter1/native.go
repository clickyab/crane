package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
)

// IsNativeAd tells if an ad is native
func IsNativeAd(c *builder.context, in models.AdData) bool {
	return in.AdType == models.NativeAdType
}
