package rtb

import (
	"github.com/clickyab/services/config"
)

var (
	develMode     = config.RegisterBoolean("devel_mode", true, "development mode?")
	adCTREffect   = config.RegisterInt("clickyab.ad_ctr_effect", 70, "")
	slotCTREffect = config.RegisterInt("clickyab.slot_ctr_effect", 30, "")
)
