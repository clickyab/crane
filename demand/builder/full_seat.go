package builder

import (
	"time"

	"clickyab.com/crane/demand/entity"
)

// SetFullSeats is a setter for seats in click and show
func SetFullSeats(pubID string, size int32, hash string, adid, cpid, cpadid int32, cpname string,
	bid float64, impTime int64, cpm, scpm float64, rt entity.RequestType, tr string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		ctr := o.publisher.CTR(size)
		ts := time.Unix(impTime, 0)
		o.seats = append(o.seats, &seat{
			context:     o,
			size:        size,
			publicID:    pubID,
			ctr:         ctr,
			adID:        adid,
			cpID:        cpid,
			cpAdID:      cpadid,
			bid:         float64(bid),
			reserveHash: hash,
			minBid:      bid,
			rate:        1,
			impTime:     ts,
			cpm:         cpm,
			scpm:        scpm,
			requestType: rt,
			cpName:      cpname,
			tr:          tr,
		})
		return o, nil
	}
}
