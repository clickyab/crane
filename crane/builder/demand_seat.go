package builder

import "fmt"

// SetDemandSeats try to add demand seat
func SetDemandSeats(pubID string, size int) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		ctr := o.publisher.CTR(size)
		if ctr <= 0 {
			return nil, fmt.Errorf("wrong ctr calculation")
		}
		ir := o.location.Country().Valid && o.location.Country().ISO == "IR"
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
		})

		return o, nil
	}
}
