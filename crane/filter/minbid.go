package filter

import (
	"math"

	"clickyab.com/crane/crane/entity"
)

// MinBID remove lesser bid
type MinBID struct {
}

// Check remove lesser bid
func (*MinBID) Check(c entity.Context, in entity.Advertise) bool {
	if c.Publisher() == nil {
		return true
	}
	t := c.Publisher().MinBid()
	if c.MinBIDPercentage() > 0 {
		t = int64(math.Ceil(float64(c.MinBIDPercentage()) / 100.0 * float64(t)))
	}

	return in.Campaign().MaxBID() >= t
}
