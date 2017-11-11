package builder

import (
	"clickyab.com/crane/crane/entity"
)

// SetFullSeats is a setter for seats in click and show
func SetFullSeats(pubID string, size int, hash string, ad entity.Advertise, bid float64) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		ctr := o.publisher.CTR(size)
		ir := o.location.Country().Valid && o.location.Country().ISO == "IR"
		showT := 0
		if o.noShowT {
			showT = 2
		}
		o.seats = append(o.seats, &seat{
			ua:        o.ua,
			parent:    o.parent,
			tid:       o.tid,
			host:      o.host,
			iran:      ir,
			alexa:     o.alexa,
			mobile:    o.os.Mobile,
			size:      size,
			publicID:  pubID,
			protocol:  o.protocol,
			ip:        o.ip,
			ref:       o.referrer,
			publisher: o.publisher,
			ftype:     o.typ,
			ctr:       ctr,

			winnerAd:    ad,
			bid:         bid,
			reserveHash: hash,
			susp:        o.suspicious,
			showT:       showT,
			// No need to following data
			minBidPercentage: o.MinBIDPercentage(),
			minBid:           bid,
			rate:             1,
		})
		return o, nil
	}
}
