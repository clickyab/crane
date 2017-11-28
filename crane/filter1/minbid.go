package filter

import (
	"math"

	"clickyab.com/crane/crane/entity"
)

//CheckMinBid check bid
func CheckMinBid(c entity.Context, in entity.Advertise) bool {
	if c.Data().Website == nil {
		return true
	}
	t := c.Data().Website.WMinBid
	if c.RTB().MinBidPercentage > 0 {
		t = int64(math.Ceil(c.RTB().MinBidPercentage * float64(t)))
	}

	return in.Campaign().MaxBID() >= t
}
