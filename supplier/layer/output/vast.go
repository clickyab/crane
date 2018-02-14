package output

import (
	"context"

	"encoding/xml"

	"fmt"
	"strings"

	"time"

	"net/http"

	"github.com/bsm/openrtb"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/xlog"
	"github.com/rs/vast"
	"github.com/rs/vmap"
)

// Seat is a video seat data
type Seat struct {
	IDExtra  string
	Start    string
	Type     string
	Duration int
	Skip     int
}

func getVast(x string) (*vast.VAST, error) {
	v := vast.VAST{}
	err := xml.Unmarshal([]byte(x), &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func stringToOffset(s string) (vmap.Offset, error) {
	if s == "start" {
		return vmap.Offset{
			Position: vmap.OffsetStart,
		}, nil
	}
	if s == "end" {
		return vmap.Offset{
			Position: vmap.OffsetEnd,
		}, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return vmap.Offset{}, err
	}
	vd := vast.Duration(d)
	return vmap.Offset{
		Duration: &vd,
	}, nil
}

// RenderVMAP is a function to handling VMAP render
func RenderVMAP(ctx context.Context, w http.ResponseWriter, resp *openrtb.BidResponse, seats map[string]Seat) error {
	v := vmap.VMAP{
		Version: "1.0",
	}

	breaks := make(map[string]*vast.VAST)

	for i := range resp.SeatBid {
		AllBid := resp.SeatBid[i].Bid
		if len(AllBid) < 1 {
			continue
		}
		bid := AllBid[0]
		markup := strings.Replace(bid.AdMarkup, "${AUCTION_PRICE}", fmt.Sprint(bid.Price), -1)
		vs, err := getVast(markup)
		if err != nil {
			xlog.GetWithError(ctx, err).Error("fail decoding vast from ad-markup")
		}
		breaks[bid.ImpID] = vs
	}

	for i := range seats {
		vs, ok := breaks[i]
		if !ok {
			continue
		}
		off, err := stringToOffset(seats[i].Start)
		if err != nil {
			continue
		}
		t := true
		ad := vmap.AdBreak{
			AdSource: &vmap.AdSource{
				ID:              i,
				VASTAdData:      vs,
				FollowRedirects: &t,
			},
			TimeOffset: off,
			BreakID:    <-random.ID,
			BreakType:  seats[i].Type,
			Extensions: &vmap.Extensions{
				Extensions: []vmap.Extension{
					{
						Type: "skip",
						Data: []byte(fmt.Sprint(seats[i].Skip)),
					},
				},
			},
		}
		v.AdBreaks = append(v.AdBreaks, ad)
	}

	w.Header().Set("content-type", "application/xml")
	b, err := xml.Marshal(v)
	assert.Nil(err)
	_, _ = w.Write(b)

	return nil
}
