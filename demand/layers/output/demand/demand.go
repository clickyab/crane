package demand

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"

	"github.com/clickyab/services/xlog"
	"github.com/golang/protobuf/jsonpb"

	"clickyab.com/crane/metrics"
	openrtb "clickyab.com/crane/openrtb/v2.5"
	"github.com/clickyab/services/version"
	"github.com/prometheus/client_golang/prometheus"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/assert"
	"github.com/rs/vast"
)

var vs = version.GetVersion()

func cdata(in string) vast.CDATAString {
	return vast.CDATAString{
		CDATA: in,
	}
}

func nativeMarkup(ctx entity.Context, s entity.NativeSeat, ver int) *openrtb.BidResponse_SeatBid {
	v := &openrtb.NativeResponse{
		Link: &openrtb.NativeResponse_Link{
			Url: s.ClickURL().String(),
		},
		Ver:         "1.1",
		Imptrackers: []string{s.ImpressionURL().String()},
	}

	for _, f := range s.Filters() {
		// TODO : Decide for duplicate assets per type :/
		a := s.WinnerAdvertise().Asset(f.Type, f.SubType, f.Extra...)
		if len(a) > 0 {
			as := &openrtb.NativeResponse_Asset{
				Id: f.ID,
				Required: func() int32 {
					if f.Required {
						return 1
					}
					return 0
				}(),
			}

			if f.Type == entity.AssetTypeImage {
				src := a[0].Data
				if ctx.Protocol() == entity.HTTPS {
					src = strings.Replace(src, "http://", "https://", -1)
				}
				as.AssetOneof = &openrtb.NativeResponse_Asset_Img{
					Img: &openrtb.NativeResponse_Asset_Image{
						Url: src,
						H:   a[0].Height,
						W:   a[0].Width,
					},
				}
			} else if f.Type == entity.AssetTypeVideo {
				// TODO : support for video VASTTAG
				as.AssetOneof = &openrtb.NativeResponse_Asset_Video_{}
			} else if f.Type == entity.AssetTypeText && f.SubType == entity.AssetTypeTextSubTypeTitle {
				as.AssetOneof = &openrtb.NativeResponse_Asset_Title_{
					Title: &openrtb.NativeResponse_Asset_Title{
						Text: a[0].Data,
						Len:  int32(len(a[0].Data)),
					},
				}
			} else if f.Type == entity.AssetTypeText {
				as.AssetOneof = &openrtb.NativeResponse_Asset_Data_{
					Data: &openrtb.NativeResponse_Asset_Data{
						Value: a[0].Data,
					},
				}
			}
			if f.Required {
				assert.True(as.AssetOneof != nil)
			}
			v.Assets = append(v.Assets, as)
		}
	}

	if ver == 0 {
		var jm = jsonpb.Marshaler{}
		st, err := jm.MarshalToString(v)
		if err != nil {
			xlog.GetWithError(context.Background(), err)
			return nil
		}
		return &openrtb.BidResponse_SeatBid{
			Bid: []*openrtb.BidResponse_SeatBid_Bid{
				{
					Id:    s.ReservedHash(),
					Impid: s.PublicID(),
					AdmOneof: &openrtb.BidResponse_SeatBid_Bid_Adm{
						Adm: st,
					},
					Adid:  fmt.Sprint(s.WinnerAdvertise().ID()),
					H:     int32(s.Height()),
					W:     int32(s.Width()),
					Price: s.CPM() / ctx.Rate(),
					Cid:   fmt.Sprint(s.WinnerAdvertise().Campaign().ID()),
					Nurl:  s.WinNoticeRequest().String(),
				},
			},
		}
	}

	return &openrtb.BidResponse_SeatBid{
		Bid: []*openrtb.BidResponse_SeatBid_Bid{
			{
				Id:    s.ReservedHash(),
				Impid: s.PublicID(),
				AdmOneof: &openrtb.BidResponse_SeatBid_Bid_AdmNative{
					AdmNative: v,
				},
				Adid:  fmt.Sprint(s.WinnerAdvertise().ID()),
				H:     int32(s.Height()),
				W:     int32(s.Width()),
				Price: s.CPM() / ctx.Rate(),
				Cid:   fmt.Sprint(s.WinnerAdvertise().Campaign().ID()),
				Nurl:  s.WinNoticeRequest().String(),
			},
		},
	}

}

func vastMarkup(ctx entity.Context, s entity.VastSeat) *openrtb.BidResponse_SeatBid {
	cv := vast.Creative{
		ID:   s.ReservedHash(),
		AdID: fmt.Sprint(s.WinnerAdvertise().ID()),
	}

	click := s.ClickURL()
	var tracking = &url.URL{}
	*tracking = *click
	q := tracking.Query()
	q.Add("tv", "1")
	src := s.WinnerAdvertise().Media()
	src = strings.Replace(src, "http://", "https://", -1)
	//
	// if ctx.Protocol() == entity.HTTPS {
	// 	src = strings.Replace(src, "http://", "https://", -1)
	// }
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
					Height:   int(s.Height()),
					Width:    int(s.Width()),
					Type:     s.WinnerAdvertise().MimeType(),
					URI:      src,
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
					URI: func() string {
						// TODO : it should check bcause we hard coded click to http
						// remove it after you remove hard coded http in click url
						return strings.Replace(tracking.String(), "http://", "https://", -1)
						//
						// if ctx.Protocol() == entity.HTTPS {
						// 	return strings.Replace(tracking.String(), "http://", "https://", -1)
						// }
						// return tracking.String()
					}(),
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
					Pricing: &vast.Pricing{
						Value: "${AUCTION_PRICE}",
					},
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
	return &openrtb.BidResponse_SeatBid{
		Bid: []*openrtb.BidResponse_SeatBid_Bid{
			{
				Id:    s.ReservedHash(),
				Impid: s.PublicID(),
				AdmOneof: &openrtb.BidResponse_SeatBid_Bid_Adm{
					Adm: string(res),
				},
				Adid:  fmt.Sprint(s.WinnerAdvertise().ID()),
				H:     s.Height(),
				W:     s.Width(),
				Price: s.CPM() / ctx.Rate(),
				Cid:   fmt.Sprint(s.WinnerAdvertise().Campaign().ID()),
				Nurl:  s.WinNoticeRequest().String(),
			},
		},
	}
}

const bannerMarkupWithoutIframe = `
	<div>
        <a href='%s' target='_blank'>
            <img src='%s'/‎>
            <img style="display: none;" src='%s'/‎>
        </a>
    </div>`

func bannerMarkup(ctx entity.Context, s entity.Seat) *openrtb.BidResponse_SeatBid {
	adMarkup := fmt.Sprintf(
		`<iframe src="%s&scpm=${AUCTION_PRICE}" width="%d" height="%d" frameborder="0"  scrolling="no" style="max-width:100%%"></iframe>`,
		s.ImpressionURL().String(),
		s.Width(),
		s.Height(),
	)
	// check for banner markup
	if ctx.BannerMarkup() {
		adMarkup = fmt.Sprintf(bannerMarkupWithoutIframe, s.ClickURL().String(), strings.Replace(s.WinnerAdvertise().Media(), "http://", "https://", -1), s.ImpressionURL().String())
	}
	return &openrtb.BidResponse_SeatBid{
		Bid: []*openrtb.BidResponse_SeatBid_Bid{
			{
				Id:    s.ReservedHash(),
				Impid: s.PublicID(),
				AdmOneof: &openrtb.BidResponse_SeatBid_Bid_Adm{
					Adm: adMarkup,
				},
				Adid:  fmt.Sprint(s.WinnerAdvertise().ID()),
				Crid:  fmt.Sprint(s.WinnerAdvertise().ID()),
				H:     s.Height(),
				W:     s.Width(),
				Price: s.CPM() / ctx.Rate(),
				Cid:   fmt.Sprint(s.WinnerAdvertise().Campaign().ID()),
			},
		},
	}

}

// Render write open-rtb bid-response to writer
func Render(_ context.Context, ctx entity.Context, rid string, ver int) (*openrtb.BidResponse, error) {
	var r []*openrtb.BidResponse_SeatBid
	for _, v := range ctx.Seats() {
		// What if we have no ad for them?
		if v.WinnerAdvertise() == nil {
			go metrics.Size.With(prometheus.Labels{
				"sup":  ctx.Publisher().Supplier().Name(),
				"size": "NaN",
				"io":   "out",
			}).Inc()
			continue
		}
		go metrics.Size.With(prometheus.Labels{
			"sup":  ctx.Publisher().Supplier().Name(),
			"size": fmt.Sprintf("%dx%d", v.Width(), v.Height()),
			"io":   "out",
		}).Inc()

		go metrics.Campaigns.With(prometheus.Labels{
			"sup": ctx.Publisher().Supplier().Name(),
			"cid": fmt.Sprint(v.WinnerAdvertise().Campaign().ID()),
		}).Inc()
		var bid *openrtb.BidResponse_SeatBid
		switch v.RequestType() {
		case entity.RequestTypeBanner:
			bid = bannerMarkup(ctx, v)
		case entity.RequestTypeVast:
			bid = vastMarkup(ctx, v.(entity.VastSeat))
		case entity.RequestTypeNative:
			bid = nativeMarkup(ctx, v.(entity.NativeSeat), ver)
		}

		if bid != nil {
			r = append(r, bid)
		}
	}

	return &openrtb.BidResponse{
		Id:      rid,
		Cur:     ctx.Currency(),
		Seatbid: r,
	}, nil
}
