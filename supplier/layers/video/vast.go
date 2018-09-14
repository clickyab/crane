package video

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	website "clickyab.com/crane/models/clickyabwebsite"
	"clickyab.com/crane/models/staticseat"
	"clickyab.com/crane/openrtb"
	"clickyab.com/crane/supplier/client"
	"clickyab.com/crane/supplier/layers/entities"
	"clickyab.com/crane/supplier/layers/internal/supplier"
	"clickyab.com/crane/supplier/layers/output"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/simplehash"
	"github.com/clickyab/services/xlog"
	"github.com/mssola/user_agent"
)

var (
	// XXX : currently, there is no need to change the endpoints per type, but if you need it, do it :) its not a rule or something.
	server = config.RegisterString("crane.supplier.banner.url", "", "route for banner")
	method = config.RegisterString("crane.supplier.banner.method", "POST", "method for banner request")
)

func writesErrorStatus(w http.ResponseWriter, status int, detail string) {
	assert.False(status == http.StatusOK)
	w.WriteHeader(status)
	_, _ = fmt.Fprint(w, detail)
}

var sup = supplier.NewClickyab()

//	d		: domain
//  a 		: public id
//	p		: current page
//	r		: ref
//	l		: length
//	tid		: tracking id
//  mimes   : comma separated accepted mime types
func vast(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pubID := r.URL.Query().Get("a")
	pub, err := website.GetWebSite(sup, pubID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	l := r.URL.Query().Get("p")
	if l == "" {
		l = r.Referer()
	}
	ref := r.URL.Query().Get("r")
	dnt, _ := strconv.Atoi(r.Header.Get("DNT"))
	tid := r.URL.Query().Get("tid")
	ln := r.URL.Query().Get("l")
	var mimes []string
	if mim := strings.Trim(r.URL.Query().Get("mimes"), "\n\t "); mim != "" {
		mimes = strings.Split(mim, ",")
	}

	type staticSeat struct {
		key        string
		seat       map[string]output.Seat
		staticSeat entities.StaticSeat
	}

	imps, seats := getImps(r, pub, getSlots(ln), mimes...)

	//<VAST version="3">
	//<Ad id="bac0d290eb20ed608ddb01ada293dcc3ae40d59d">
	//<InLine>
	//<AdSystem version="115"><![CDATA[3rdAd]]></AdSystem>
	//<AdTitle><![CDATA[Test Campaign22]]></AdTitle>
	//<Impression />
	//<Creatives>
	//<Creative id="5ae5c712d0c45" AdID="1ae5c712d0c25">
	//<Linear skipoffset="00:00:03">
	//<Duration>00:00:28</Duration>
	//<TrackingEvents>
	//<Tracking event="complete" />
	//</TrackingEvents>
	//<VideoClicks>
	//<ClickThrough><![CDATA[https://goo.gl/Ur8w1Z]]></ClickThrough>
	//</VideoClicks>
	//<MediaFiles>
	//<MediaFile delivery="streaming" type="image/jpeg" width="800" height="440"><![CDATA[http://static.clickyab.com/ad/800x440/20180704-6915843-0934823740273336e888788e798d6dc0f8bded34.jpg]]></MediaFile>
	//</MediaFiles>
	//</Linear>
	//<CreativeExtensions />
	//</Creative>
	//</Creatives>
	//<Description />
	//<Survey />
	//<Pricing>10</Pricing>
	//<Extensions />
	//</InLine>
	//</Ad>
	//</VAST>

	var finalStaticSeats = make([]staticSeat, 0)

	for j := range seats {
		d, exists := staticseat.GetMultiStaticSeat(pub, "vast", seats[j].Start)
		if exists {
			if d[0].Fix() { // at least one is fix we should exactly return 1 from those no matter what the chance is
				finalStaticSeats = append(finalStaticSeats, staticSeat{
					staticSeat: alwaysReturnFix(d),
					seat: map[string]output.Seat{
						j: {
							Type:     seats[j].Type,
							Duration: seats[j].Duration,
							IDExtra:  seats[j].IDExtra,
							Skip:     seats[j].Skip,
							Start:    seats[j].Start,
						},
					},
					key: j,
				})
				delete(seats, j)
			} else {
				if d[0].Chance() > rand.Intn(100) {
					finalStaticSeats = append(finalStaticSeats, staticSeat{
						staticSeat: d[0],
						seat: map[string]output.Seat{
							j: {
								Type:     seats[j].Type,
								Duration: seats[j].Duration,
								IDExtra:  seats[j].IDExtra,
								Skip:     seats[j].Skip,
								Start:    seats[j].Start,
							},
						},
						key: j,
					})
					delete(seats, j)
				}
			}
		}
	}

	ua := user_agent.New(r.UserAgent())

	rIP := framework.RealIP(r)
	rUserAgent := r.UserAgent()

	bq := &openrtb.BidRequest{
		Id: <-random.ID,
		User: &openrtb.User{
			Id: vastUserIDGenerator(tid, rUserAgent, rIP),
		},
		Imp: imps,
		DistributionchannelXoneof: &openrtb.BidRequest_Site{

			Site: &openrtb.Site{
				Mobile: func() int32 {
					if ua.Mobile() {
						return 1
					}
					return 0
				}(),
				Page:   l,
				Ref:    ref,
				Domain: pub.Name(),
				Name:   pub.Name(),
				Id:     fmt.Sprint(pub.ID()),
				Cat:    pub.Categories(),
			},
		},
		Device: &openrtb.Device{
			Ip:  rIP,
			Dnt: int32(dnt),
			Os:  ua.OS(),
			Ua:  rUserAgent,
		},
		Ext: &openrtb.BidRequest_Ext{
			Capping:    openrtb.Capping_Reset,
			Underfloor: true,
		},
	}

	var br = &openrtb.BidResponse{}

	br, err = client.Call(ctx, method.String(), server.String(), bq)
	if err != nil {
		if len(finalStaticSeats) > 0 {
			br = &openrtb.BidResponse{}
		} else {
			e := "demand error"
			writesErrorStatus(w, http.StatusInternalServerError, e)
			xlog.GetWithError(ctx, err).Debugf(e)
			return
		}
	}

	for i := range finalStaticSeats {
		seats[finalStaticSeats[i].key] = finalStaticSeats[i].seat[finalStaticSeats[i].key]
		br.Seatbid = append(br.Seatbid, &openrtb.BidResponse_SeatBid{
			Bid: []*openrtb.BidResponse_SeatBid_Bid{
				{
					Id:    <-random.ID,
					Impid: finalStaticSeats[i].key,
					AdmXoneof: &openrtb.BidResponse_SeatBid_Bid_Adm{
						Adm: finalStaticSeats[i].staticSeat.RTBMarkup(),
					},
				},
			},
		})
	}

	if err := output.RenderVMAP(ctx, w, br, seats); err != nil {
		e := "render failed"
		writesErrorStatus(w, http.StatusInternalServerError, e)
		xlog.GetWithError(ctx, err).Debugf(e)
		return
	}
}

func alwaysReturnFix(seats []entities.StaticSeat) entities.StaticSeat {
	n := rand.Int() % len(seats)
	return seats[n]
}

// vastUserIDGenerator create user id for vast
func vastUserIDGenerator(tid, ua, ip string) string {
	return simplehash.MD5(fmt.Sprintf("%s%s%s", tid, ua, ip))
}
