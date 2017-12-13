package builder

import (
	"fmt"

	"clickyab.com/crane/crane/builder/internal/cyslot"
)

// SetDemandSeats try to add demand seat
func SetDemandSeats(sd map[string]string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		ir := o.location.Country().Valid && o.location.Country().ISO == "IR"
		for i := range sd {
			size, err := cyslot.GetSize(sd[i])
			if err != nil {
				return nil, err
			}

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
				publicID:  i,
				protocol:  o.protocol,
				ip:        o.ip,
				ref:       o.referrer,
				publisher: o.publisher,
				ftype:     o.typ,
				ctr:       o.publisher.CTR(size),
				showT:     showT,
				rate:      o.rate,
			})
		}
		if len(o.seats) == 0 {
			return nil, fmt.Errorf("no supported seat")
		}
		return o, nil
	}
}
