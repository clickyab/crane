package filter

import (
	"math"

	"clickyab.com/crane/crane/builder"
	"clickyab.com/gad/models"
)

//CheckMinBid find isp
func CheckMinBid(c *builder.context, in models.AdData) bool {
	if c.GetData().Website == nil {
		return true
	}
	t := c.GetData().Website.WMinBid
	if c.GetRTB().MinBidPercentage > 0 {
		t = int64(math.Ceil(c.GetRTB().MinBidPercentage * float64(t)))
	}

	return in.CampaignMaxBid >= t
}
