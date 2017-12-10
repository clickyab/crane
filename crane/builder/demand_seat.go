package builder

import (
	"errors"
	"fmt"

	"clickyab.com/crane/crane/builder/internal/cyslot"
)

// SeatDetail carries data about a seat
type SeatDetail struct {
	PubID string
	Size  int
	W     int
	H     int
}

// SetDemandSeats try to add demand seat
func SetDemandSeats(sd []SeatDetail) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		ir := o.location.Country().Valid && o.location.Country().ISO == "IR"
		for i := range sd {
			if sd[i].Size == 0 {
				if o.typ == "native" {
					sd[i].Size = 20
				} else {
					var err error
					sd[i].Size, err = cyslot.GetSize(fmt.Sprintf("%dx%d", sd[i].W, sd[i].H))
					if err != nil {
						return nil, errors.New("invalid size requested")
					}
				}
			}

			ctr := o.publisher.CTR(sd[i].Size)
			if ctr <= 0 {
				return nil, fmt.Errorf("wrong ctr calculation")
			}

			o.seats = append(o.seats, &seat{
				ua:              o.ua,
				parent:          o.parent,
				tid:             o.tid,
				host:            o.host,
				iran:            ir,
				alexa:           o.alexa,
				mobile:          o.os.Mobile,
				size:            sd[i].Size,
				publicID:        sd[i].PubID,
				protocol:        o.protocol,
				ip:              o.ip,
				publisherDomain: o.publisher.Name(),
				ref:             o.referrer,
				supplier:        o.publisher.Supplier(),
				ftype:           o.typ,
				ctr:             o.publisher.CTR(sd[i].Size),
			})
		}
		return o, nil
	}
}
