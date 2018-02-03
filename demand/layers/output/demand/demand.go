package demand

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"encoding/xml"

	"net/url"

	"clickyab.com/crane/demand/entity"
	"github.com/bsm/openrtb"
	"github.com/bsm/openrtb/native/response"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/version"
	"github.com/rs/vast"
)

var vs = version.GetVersion()

func cdata(in string) vast.CDATAString {
	return vast.CDATAString{
		CDATA: in,
	}
}

func nativeMarkup(ctx entity.Context, s entity.NativeSeat) *openrtb.Bid {
	v := response.Response{
		Link: response.Link{
			URL: s.ClickURL().String(),
		},
		ImpTrackers: []string{s.ImpressionURL().String()},
		Ver:         "1.1",
	}
	for _, f := range s.Filters() {
		// TODO : Decide for duplicate assets per type :/
		a := s.WinnerAdvertise().Assets(f.Type, f.SubType, f.Extra...)
		if len(a) > 0 {
			req := 0
			if f.Required {
				req = 1
			}
			as := response.Asset{
				ID:       f.ID,
				Required: req,
			}
			if f.Type == entity.AssetTypeImage {
				as.Image = &response.Image{
					URL:    a[0].Data,
					Height: a[0].Height,
					Width:  a[0].Width,
				}
			} else if f.Type == entity.AssetTypeVideo {
				// TODO : support for video VASTTAG
				as.Video = &response.Video{}
			} else if f.Type == entity.AssetTypeText && f.SubType == entity.AssetTypeTextSubTypeTitle {
				as.Title = &response.Title{
					Text: a[0].Data,
				}
			} else if f.Type == entity.AssetTypeText {
				as.Data = &response.Data{
					Value: a[0].Data,
				}
			}

			v.Assets = append(v.Assets, as)
		}
	}

	res, err := json.Marshal(v)
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
		NURL:       s.WinNoticeRequest().String(),
	}
}

func vastMarkup(ctx entity.Context, s entity.VastSeat) *openrtb.Bid {
	cv := vast.Creative{
		ID:   s.ReservedHash(),
		AdID: fmt.Sprint(s.WinnerAdvertise().ID()),
	}

	click := s.ClickURL()
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
					Height:   s.Height(),
					Width:    s.Width(),
					Type:     s.WinnerAdvertise().MimeType(),
					URI:      s.WinnerAdvertise().Media(),
					Delivery: "streaming",
				},
			},
			VideoClicks: &vast.VideoClicks{
				ClickThroughs: []vast.VideoClick{
					{URI: s.ClickURL().String()},
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
					Impressions: []vast.Impression{
						{URI: s.ImpressionURL().String()},
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
		NURL:       s.WinNoticeRequest().String(),
	}
}

func bannerMarkup(ctx entity.Context, s entity.Seat) *openrtb.Bid {
	return &openrtb.Bid{
		ID:    s.ReservedHash(),
		ImpID: s.PublicID(),
		AdMarkup: fmt.Sprintf(
			`<iframe src="%s&scpm=${AUCTION_PRICE}" width="%d" height="%d" frameborder="0"  scrolling="no"></iframe>`,
			s.ImpressionURL().String(),
			s.Width(),
			s.Height(),
		),
		AdID:       fmt.Sprint(s.WinnerAdvertise().ID()),
		CreativeID: fmt.Sprint(s.WinnerAdvertise().ID()),
		H:          s.Height(),
		W:          s.Width(),
		Price:      s.CPM() / ctx.Rate(),
		CampaignID: openrtb.StringOrNumber(fmt.Sprint(s.WinnerAdvertise().Campaign().ID())),
	}
}

// Render write open-rtb bid-response to writer
func Render(_ context.Context, w http.ResponseWriter, ctx entity.Context) error {
	var r []openrtb.SeatBid
	w.Header().Set("crane-version", fmt.Sprint(vs.Count))
	for _, v := range ctx.Seats() {
		// What if we have no ad for them?
		if v.WinnerAdvertise() == nil {
			continue
		}
		var bid *openrtb.Bid
		switch v.SubType() {
		case entity.RequestTypeWeb, entity.RequestTypeApp:
			bid = bannerMarkup(ctx, v)
		case entity.RequestTypeVast:
			bid = vastMarkup(ctx, v.(entity.VastSeat))
		case entity.RequestTypeNative:
			bid = nativeMarkup(ctx, v.(entity.NativeSeat))
		}

		if bid != nil {
			r = append(r, openrtb.SeatBid{Bid: []openrtb.Bid{*bid}})
		}
	}
	w.Header().Set("content-type", "application/json")
	j := json.NewEncoder(w)
	return j.Encode(openrtb.BidResponse{
		Currency: ctx.Currency(),
		ID:       <-random.ID,
		SeatBid:  r,
	})
}
