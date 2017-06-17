package rtb

import (
	"time"

	"github.com/clickyab/services/config"
)

var (
	megaImpressionTTL = config.RegisterDuration("clickyab.rtb.mega_imp_ttl", 72*time.Hour, "clickyab rtb mega impression time to live")
	adCtrEffect       = config.RegisterInt("clickyab.rtb.ad_ctr_effect", 30, "clicktab rtb advertise controller effect")
	slotCtrEffect     = config.RegisterInt("clickyab.rtb.slot_ctr_effect", 70, "clickyab rtb slot controller effect")
)
