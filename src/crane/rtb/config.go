package rtb

import (
	"services/config"
	"time"
)

var (
	megaImpressionTTL = config.RegisterDuration("clickyab.rtb.mega_imp_ttl", 72*time.Hour)
	adCtrEffect       = config.RegisterInt("clickyab.rtb.ad_ctr_effect", 30)
	slotCtrEffect     = config.RegisterInt("clickyab.rtb.slot_ctr_effect", 70)
)
