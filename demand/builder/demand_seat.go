package builder

import (
	"fmt"
	"time"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/internal/cyslot"
	"github.com/bsm/openrtb"
	"github.com/bsm/openrtb/native/request"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/kv"
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
	MinCPC float64
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

func incShare(sup entity.Supplier, price float64) float64 {
	return (price * float64(sup.Share())) / 100
}

func decShare(sup entity.Supplier, price float64) float64 {
	return (price * 100.0) / float64(sup.Share())
}

// SetSeats try to add demand seat
func SetSeats(sd ...entity.Seat) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.seats = sd
		return o, nil
	}
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
				iqs := kv.NewAEAVStore(fmt.Sprintf("DEQS_%s", time.Now().Truncate(time.Hour*24).Format("060102")), time.Hour*72)
				iqs.IncSubKey(fmt.Sprintf("%s_%s_%s", o.Publisher().Supplier().Name(), time.Now().Truncate(time.Hour).Format("15"), sd[i].Size), 1)
				switch sd[i].Type {
				case entity.RequestTypeNative:
					size = 20
				case entity.RequestTypeVast:
					size = 9
				default:
					return nil, err
				}
			}

			minCPC := sd[i].MinCPC
			if minCPC == 0 {
				minCPC = float64(o.publisher.Supplier().SoftFloorCPC(fmt.Sprint(sd[i].Type), fmt.Sprint(o.publisher.Type())))
			}
			seat := seat{
				context:     o,
				size:        size,
				publicID:    sd[i].PubID,
				minBid:      sd[i].MinBid,
				ctr:         o.publisher.CTR(size),
				rate:        o.rate,
				requestType: sd[i].Type,
				minCPC:      minCPC,
				softCPM:     float64(o.Publisher().Supplier().SoftFloorCPM(fmt.Sprint(sd[i].Type), fmt.Sprint(o.Publisher().Type()))),
				minCPM:      float64(incShare(o.Publisher().Supplier(), sd[i].MinBid)),
			}

			if seat.softCPM < seat.minCPM {
				seat.softCPM = seat.minCPM
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
