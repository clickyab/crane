package builder

import (
	"fmt"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/internal/cyslot"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/config"
)

var (
	vastLinearDefaultLen  = config.RegisterInt("crane.vast.linear_len", 30, "default len of linear ad in sec, if not presented in request")
	vastLinearDefaultSkip = config.RegisterInt("crane.vast.linear_skip", 10, "default len of linear ad in sec, if not presented in request")
	//vastNonLinearDefaultLen = config.RegisterInt("crane.vast.non_linear_len", 10, "default len of non linear ad in sec, if not presented in request")
)

// SeatType means the seat is banner or video
type SeatType int

const (
	// SeatTypeBanner means this is a banner
	SeatTypeBanner SeatType = iota
	// SeatTypeVideo means this seat is video
	SeatTypeVideo
)

// DemandSeatData is a struct needed for create a demand seat
type DemandSeatData struct {
	PubID  string
	Size   string
	MinBid float64
	Type   SeatType
	Video  *openrtb.Video
	Banner *openrtb.Banner
}

func coalesce(v ...int) int {
	for i := range v {
		if v[i] > 0 {
			return v[i]
		}
	}
	return 0
}

// SetDemandSeats try to add demand seat
func SetDemandSeats(sd ...DemandSeatData) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		ir := o.location.Country().Valid && o.location.Country().ISO == "IR"
		for i := range sd {
			var (
				size int
				err  error
			)
			linear := false
			if sd[i].Type == SeatTypeVideo && sd[i].Video.Linearity == 1 {
				linear = true
			}
			size, err = cyslot.GetSize(sd[i].Size)
			if err != nil {
				if !linear {
					return nil, err
				}
				size = 9
			}

			showT := 0
			if o.noShowT {
				showT = 2
			}

			seat := seat{
				ua:               o.ua,
				parent:           o.parent,
				tid:              o.tid,
				host:             o.host,
				iran:             ir,
				alexa:            o.alexa,
				mobile:           o.os.Mobile,
				size:             size,
				publicID:         sd[i].PubID,
				minBid:           sd[i].MinBid,
				protocol:         o.protocol,
				ip:               o.ip,
				ref:              o.referrer,
				publisher:        o.publisher,
				ftype:            o.typ,
				subType:          o.subTyp,
				ctr:              o.publisher.CTR(size),
				showT:            showT,
				rate:             o.rate,
				minBidPercentage: o.MinBIDPercentage(),
				fatFinger:        o.fatFinger,
			}
			if sd[i].Type == SeatTypeVideo {
				seat.subType = entity.RequestTypeVast
				seat.mimes = sd[i].Video.Mimes
				o.seats = append(o.seats, &vastSeat{
					seat:      seat,
					linear:    linear,
					duration:  vastLinearDefaultLen.Int(),
					skipAfter: coalesce(sd[i].Video.SkipMin, vastLinearDefaultSkip.Int()),
				})
			} else {
				o.seats = append(o.seats, &seat)
			}
		}
		if len(o.seats) == 0 {
			return nil, fmt.Errorf("no supported seat")
		}
		return o, nil
	}
}
