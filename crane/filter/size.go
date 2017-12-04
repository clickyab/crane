package filter

import (
	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
	"clickyab.com/gad/utils"
)

// CheckWebSize check if the banner size exists in the request
func CheckWebSize(c *builder.context, in models.AdData) bool {
	if in.AdType == utils.AdTypeVideo {
		for _, size := range c.GetRTB().Slots {
			if utils.InVideoSize(size.Size) {
				return true
			}
		}
		return false
	}

	for _, size := range c.GetRTB().Slots {
		if size.Size == in.AdSize {
			return true
		}
	}
	return false
}

// CheckVastSize check if the banner size fits for Vast Template
func CheckVastSize(_ *builder.context, in models.AdData) bool {
	if in.AdType == utils.AdTypeDynamic {
		return false
	}

	return in.AdType == utils.AdTypeVideo || utils.InVastSize(in.AdSize)
}

// CheckAppSize check if the banner size exists in the request
func CheckAppSize(c *builder.context, in models.AdData) bool {
	if in.AdType == utils.AdTypeVideo || in.AdType == utils.AdTypeDynamic {
		return false
	}

	for _, size := range c.GetRTB().Slots {
		if size.Size == in.AdSize {
			return true
		}
	}
	return false
}

// CheckWebMobileSize check if the banner size exists in the request
func CheckWebMobileSize(c *builder.context, in models.AdData) bool {
	return in.AdSize == 8
}
