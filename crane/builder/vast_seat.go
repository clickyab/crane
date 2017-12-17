package builder

import (
	"fmt"

	"clickyab.com/crane/crane/builder/internal/cyvast"
)

// SetVastSeats try to add vast seats
func SetVastSeats(l, basePubID string, first, mid, last bool) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		var pubIdSize = make(map[string]int)
		_, length := cyvast.MakeVastLen(l, first, mid, last)
		var i int
		ir := o.location.Country().Valid && o.location.Country().ISO == "IR"

		for m := range length {
			i++
			lenType := length[m][0]
			if lenType != "linear" && o.IsMobile() {
				continue
			}
			pub := fmt.Sprintf("%s%s", basePubID, length[m][1])

			size := cyvast.VastNonLinearSize
			//stop := ""
			if lenType == "linear" {
				size = cyvast.VastLinearSize
				//	stop = length[m][3]
			}
			pubIdSize[pub] = size

			o.seats = append(o.seats, &seat{
				ua:               o.ua,
				parent:           o.parent,
				tid:              o.tid,
				host:             o.host,
				iran:             ir,
				alexa:            o.alexa,
				mobile:           o.os.Mobile,
				size:             size,
				publicID:         pub,
				minBid:           float64(o.publisher.MinBid()),
				protocol:         o.protocol,
				ip:               o.ip,
				ref:              o.referrer,
				publisher:        o.publisher,
				ftype:            o.typ,
				ctr:              o.publisher.CTR(size),
				showT:            2, // No :)
				rate:             o.rate,
				minBidPercentage: o.MinBIDPercentage(),

				showExtraParam: map[string]string{
					"pos":  m,
					"type": length[m][2],
					"l":    lenType,
				},
			})
		}

		if len(o.seats) == 0 {
			return nil, fmt.Errorf("no supported seat")
		}
		return o, nil
	}
}
