package demand

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"encoding/xml"

	"net/url"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/gad/src/version"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/random"
	"github.com/rs/vast"
)

var vs = version.GetVersion()

func cdata(in string) vast.CDATAString {
	return vast.CDATAString{
		CDATA: in,
	}
}

func vastMarkup(ctx entity.Context, s entity.VastSeat) *openrtb.Bid {
	cv := vast.Creative{
		ID:   s.ReservedHash(),
		AdID: fmt.Sprint(s.WinnerAdvertise().ID()),
	}

	click, err := url.Parse(s.ClickURL())
	assert.Nil(err)
	var tracking = &url.URL{}
	*tracking = *click
	q := tracking.Query()
	q.Add("tv", "1")
	tracking.RawQuery = q.Encode()
	if s.Linear() {
		skipAfter := vast.Duration(s.SkipAfter())
		cv.Linear = &vast.Linear{
			Duration: vast.Duration(s.Duration()),
			SkipOffset: &vast.Offset{
				Duration: &skipAfter,
			},

			MediaFiles: []vast.MediaFile{
				{
					Height: s.Height(),
					Width:  s.Width(),
					Type:   s.WinnerAdvertise().MimeType(),
					URI:    click.String(),
				},
			},

			TrackingEvents: []vast.Tracking{
				{
					URI:   tracking.String(),
					Event: "complete",
				},
			},
		}
	} else {
		return nil
	}
	v := vast.VAST{
		Version: "3",
		Ads: []vast.Ad{
			{
				ID:       s.ReservedHash(),
				Sequence: 0, // Currently we use one add per slot system
				InLine: &vast.InLine{
					AdSystem: &vast.AdSystem{
						Version: fmt.Sprint(vs.Count),
						Name:    "3rdAd",
					},
					AdTitle: cdata(s.WinnerAdvertise().Campaign().Name()),
					Pricing: "${AUCTION_PRICE}",
					Creatives: []vast.Creative{
						cv,
					},
				},
			},
		},
	}

	res, err := xml.MarshalIndent(v, "", "  ")
	assert.Nil(err)
	return &openrtb.Bid{
		ID:         s.ReservedHash(),
		ImpID:      s.PublicID(),
		AdMarkup:   string(res),
		AdID:       fmt.Sprint(s.WinnerAdvertise().ID()),
		H:          s.Height(),
		W:          s.Width(),
		Price:      s.CPM() / ctx.Rate(),
		CampaignID: openrtb.StringOrNumber(fmt.Sprint(s.WinnerAdvertise().Campaign().ID())),
	}
}

func bannerMarkup(ctx entity.Context, s entity.Seat) *openrtb.Bid {
	return &openrtb.Bid{
		ID:    s.ReservedHash(),
		ImpID: s.PublicID(),
		AdMarkup: fmt.Sprintf(
			`<iframe src="%s&scpm=${AUCTION_PRICE}" width="%d" height="%d" frameborder="0"  scrolling="no"></iframe>`,
			s.ShowURL(),
			s.Width(),
			s.Height(),
		),
		AdID:       fmt.Sprint(s.WinnerAdvertise().ID()),
		H:          s.Height(),
		W:          s.Width(),
		Price:      s.CPM() / ctx.Rate(),
		CampaignID: openrtb.StringOrNumber(fmt.Sprint(s.WinnerAdvertise().Campaign().ID())),
	}
}

// Render write open-rtb bid-response to writer
func Render(_ context.Context, w http.ResponseWriter, ctx entity.Context) error {
	r := openrtb.SeatBid{}
	w.Header().Set("crane-version", fmt.Sprint(vs.Count))
	for _, v := range ctx.Seats() {
		// What if we have no ad for them?
		if v.WinnerAdvertise() == nil {
			continue
		}
		var bid *openrtb.Bid
		switch entity.RequestType(v.SubType()) {
		case entity.RequestTypeWeb, entity.RequestTypeApp:
			bid = bannerMarkup(ctx, v)
		case entity.RequestTypeVast:
			bid = vastMarkup(ctx, v.(entity.VastSeat))
		}

		if bid != nil {
			r.Bid = append(r.Bid, *bid)
		}
	}
	w.Header().Set("content-type", "application/json")
	j := json.NewEncoder(w)
	return j.Encode(openrtb.BidResponse{
		Currency: ctx.Currency(),
		ID:       <-random.ID,
		SeatBid:  []openrtb.SeatBid{r},
	})
}
