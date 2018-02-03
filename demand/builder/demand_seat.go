package builder

import (
	"fmt"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/internal/cyslot"
	"github.com/bsm/openrtb"
	"github.com/bsm/openrtb/native/request"
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
	// SeatTypeNative means this is a native request
	SeatTypeNative
)

// DemandSeatData is a struct needed for create a demand seat
type DemandSeatData struct {
	PubID  string
	Size   string
	MinBid float64
	Type   SeatType
	Video  *openrtb.Video
	Banner *openrtb.Banner
	Assets []request.Asset
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

			seat := seat{
				context:  o,
				subType:  o.subTyp,
				size:     size,
				publicID: sd[i].PubID,
				minBid:   sd[i].MinBid,
				ctr:      o.publisher.CTR(size),
				rate:     o.rate,
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
			} else if sd[i].Type == SeatTypeNative {
				seat.subType = entity.RequestTypeNative
				o.seats = append(o.seats, &nativeSeat{
					seat:    seat,
					filters: assetToFilterFunc(sd[i].Assets),
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
