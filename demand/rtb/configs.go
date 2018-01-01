package rtb

import (
	"github.com/clickyab/services/config"
)

var (
	//develMode     = config.RegisterBoolean("devel_mode", true, "development mode?")
	adCTREffect   = config.RegisterInt("crane.rtb.ad_ctr_effect", 50, "ad ctr effect")
	slotCTREffect = config.RegisterInt("crane.rtb.slot_ctr_effect", 50, "slot ctr effect")
)
