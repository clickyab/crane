package builder

import (
	"time"

	"clickyab.com/crane/demand/entity"
)

// SetFullSeats is a setter for seats in click and show
func SetFullSeats(pubID string, size int, hash string, ad entity.Creative, bid float64, impTime int64, cpm, scpm float64) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		ctr := o.publisher.CTR(size)
		ts := time.Unix(impTime, 0)
		o.seats = append(o.seats, &seat{
			context:     o,
			subType:     o.subTyp,
			size:        size,
			publicID:    pubID,
			ctr:         ctr,
			winnerAd:    ad,
			bid:         bid,
			reserveHash: hash,
			minBid:      bid,
			rate:        1,
			impTime:     ts,
			cpm:         cpm,
			scpm:        scpm,
		})
		return o, nil
	}
}
