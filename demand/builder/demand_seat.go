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
	vastLinearDefaultSkip = config.RegisterInt("crane.vast.linear_skip", 10, "default len of linear ad in sec, if not presented in request")
	//vastNonLinearDefaultLen = config.RegisterInt("crane.vast.non_linear_len", 10, "default len of non linear ad in sec, if not presented in request")
)

// DemandSeatData is a struct needed for create a demand seat
type DemandSeatData struct {
	PubID  string
	Size   string
	MinBid float64
	Type   entity.RequestType
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
			if sd[i].Type == entity.RequestTypeVast && sd[i].Video.Linearity == 1 {
				linear = true
			}
			size, err = cyslot.GetSize(sd[i].Size)
			if err != nil {
				switch sd[i].Type {
				case entity.RequestTypeNative:
					size = 20
				case entity.RequestTypeVast:
					size = 9
				default:
					return nil, err
				}
			}

			seat := seat{
				context:     o,
				size:        size,
				publicID:    sd[i].PubID,
				minBid:      sd[i].MinBid,
				ctr:         o.publisher.CTR(size),
				rate:        o.rate,
				requestType: sd[i].Type,
			}
			if sd[i].Type == entity.RequestTypeVast {
				seat.mimes = sd[i].Video.Mimes

				// duration is field in set winner advertise
				o.seats = append(o.seats, &vastSeat{
					seat:      seat,
					linear:    linear,
					skipAfter: coalesce(sd[i].Video.SkipMin, vastLinearDefaultSkip.Int()),
				})
			} else if sd[i].Type == entity.RequestTypeNative {
				seat.requestType = entity.RequestTypeNative
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
