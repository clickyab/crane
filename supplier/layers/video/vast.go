package video

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"strings"

	"math/rand"

	website "clickyab.com/crane/models/clickyabwebsite"
	"clickyab.com/crane/models/staticseat"
	"clickyab.com/crane/supplier/client"
	"clickyab.com/crane/supplier/layers/entities"
	"clickyab.com/crane/supplier/layers/internal/supplier"
	"clickyab.com/crane/supplier/layers/output"
	"github.com/bsm/openrtb"
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
	fmt.Fprint(w, detail)
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

	//vast rtb mark up example :
	//<VAST version="3">
	//	<Ad id="awdawdawdawd">
	//		<InLine>
	//			<AdSystem version="115"><![CDATA[3rdAd]]></AdSystem>
	//			<AdTitle><![CDATA[Test Campaign22]]></AdTitle>
	//			<Impression><![CDATA[http://demand.clickyab.com/api/pixel/9b7eaf834c3bcee6907f98f466218bdd1dc0dfc0/9/demand/vast/eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhaWQiOiIxMDEiLCJiaWQiOiIyOTk5Ljk5OTk5OTk5OTk5OTUiLCJjbW9kZSI6IjEiLCJjcG0iOiIxMjM0Ni45NzAzNjcwOTQyMDQiLCJkb20iOiJqYWJlaC5jb20iLCJleHAiOiIxODA0MjkxMzI3MDgiLCJmZiI6IkYiLCJpYXQiOjE1MjQ5ODg2MjgsIm5vdyI6IjE1MjQ5ODg2MjgiLCJwaWQiOiIxNzY0NTExIiwicHQiOiJ3ZWIiLCJzdXAiOiJjbGlja3lhYiIsInN1c3AiOiIwIiwidCI6IlQiLCJ1YWlwIjoiMjdlNWQ1OGI0OWE5YjBiOTMwNWQyNTQ2OTAzODg1YTcifQ.na5Aj_x3QO3PQtTh7V-Cr-sL82kMIh1S_rjWGN2tBGbBH7DXZjJnE4j4MhbOhW3fWwYIjEwLXGGElg_boMj18bQjL0Vw7y3230_HEnjdrJ7E7rkfX-lKiRjA5z6Xp9Mirt8nMWkRKSSB0-RH3kKyoOU5djP57VMBFbpU-VOVrCXArIVQF9EvCcvI5cxvkOImIxWRDah7fTipvHEZ0OvQ17S7s2jwCKMG02o5YIobFDTNffr0TDty8oA-CM1CTSgmGY5V4K-dsbljqOpJIwkV2Y2SXmsLe5aUbA5KQLrykDkb4Gssk5DevfR5XQLjxBJmLIK0UHqZajmME0CCrKZ6_A?parent=&ref=&reg=fr&tid=0f894c5f7c]]></Impression>
	//			<Creatives>
	//			<Creative id="9b7eaf834c3bcee6907f98f466218bdd1dc0dfc0" AdID="101">
	//			<Linear skipoffset="00:00:03">
	//			<Duration>00:00:18</Duration>
	//			<TrackingEvents>
	//			<Tracking event="complete"><![CDATA[http://demand.clickyab.com/api/click/9b7eaf834c3bcee6907f98f466218bdd1dc0dfc0/9/demand/vast/eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhaWQiOiIxMDEiLCJiaWQiOiIyOTk5Ljk5OTk5OTk5OTk5OTUiLCJjbW9kZSI6IjEiLCJjcG0iOiIxMjM0Ni45NzAzNjcwOTQyMDQiLCJkb20iOiJqYWJlaC5jb20iLCJleHAiOiIxODA1MDIxMjI3MDgiLCJmZiI6IkYiLCJpYXQiOjE1MjQ5ODg2MjgsIm5vdyI6IjE1MjQ5ODg2MjgiLCJwaWQiOiIxNzY0NTExIiwicHQiOiJ3ZWIiLCJzdXAiOiJjbGlja3lhYiIsInN1c3AiOiIwIiwidCI6IlQiLCJ1YWlwIjoiMjdlNWQ1OGI0OWE5YjBiOTMwNWQyNTQ2OTAzODg1YTcifQ.Kh2uTKC1F5EznpqDenWD278XRaZEQ0YbgrfC9wF3rfYghbeLFcZsniGtDoCN9Jg5HR4Ouo0tXPLxuSiTCw5don31OAWIU2sNY_KEoFiYBkGp4Z6S-uPqQfNuwryp7blEtpTtdKAk63bpYFvxtyuL2qF8OH0tkzkbHhSWhTOD6cn8P086xVWSX3RNp-KrEycTw5bfxCWRoYfOyrBu-YXSX6Ry7eEVtCKpPJ70Da3Lrm74cR_iqDWiBvNt0V9yM-NqIpwHP3JBOwLYV46CzlaloymAN5ayavzYJhT5X7vWTQeNrC62u8crbEIS41oJ0YQdsNyPPOMcwJ7NEeaexIqU7Q?parent=&ref=&reg=fr&tid=0f894c5f7c&tv=1]]></Tracking>
	//			</TrackingEvents>
	//			<VideoClicks>
	//			<ClickThrough><![CDATA[http://demand.clickyab.com/api/click/9b7eaf834c3bcee6907f98f466218bdd1dc0dfc0/9/demand/vast/eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhaWQiOiIxMDEiLCJiaWQiOiIyOTk5Ljk5OTk5OTk5OTk5OTUiLCJjbW9kZSI6IjEiLCJjcG0iOiIxMjM0Ni45NzAzNjcwOTQyMDQiLCJkb20iOiJqYWJlaC5jb20iLCJleHAiOiIxODA1MDIxMjI3MDgiLCJmZiI6IkYiLCJpYXQiOjE1MjQ5ODg2MjgsIm5vdyI6IjE1MjQ5ODg2MjgiLCJwaWQiOiIxNzY0NTExIiwicHQiOiJ3ZWIiLCJzdXAiOiJjbGlja3lhYiIsInN1c3AiOiIwIiwidCI6IlQiLCJ1YWlwIjoiMjdlNWQ1OGI0OWE5YjBiOTMwNWQyNTQ2OTAzODg1YTcifQ.Kh2uTKC1F5EznpqDenWD278XRaZEQ0YbgrfC9wF3rfYghbeLFcZsniGtDoCN9Jg5HR4Ouo0tXPLxuSiTCw5don31OAWIU2sNY_KEoFiYBkGp4Z6S-uPqQfNuwryp7blEtpTtdKAk63bpYFvxtyuL2qF8OH0tkzkbHhSWhTOD6cn8P086xVWSX3RNp-KrEycTw5bfxCWRoYfOyrBu-YXSX6Ry7eEVtCKpPJ70Da3Lrm74cR_iqDWiBvNt0V9yM-NqIpwHP3JBOwLYV46CzlaloymAN5ayavzYJhT5X7vWTQeNrC62u8crbEIS41oJ0YQdsNyPPOMcwJ7NEeaexIqU7Q?parent=&ref=&reg=fr&tid=0f894c5f7c]]></ClickThrough>
	//			</VideoClicks>
	//			<MediaFiles>
	//			<MediaFile delivery="streaming" type="video/mp4" width="800" height="440"><![CDATA[http://static.clickyab.com/ad/video/20180110-7476975-6fc73384b929682ad3afb164bed22c9a15463881.cy]]></MediaFile>
	//			</MediaFiles>
	//			</Linear>
	//			<CreativeExtensions></CreativeExtensions>
	//			</Creative>
	//			</Creatives>
	//			<Description></Description>
	//			<Survey></Survey>
	//			<Pricing>${AUCTION_PRICE}</Pricing>
	//			<Extensions></Extensions>
	//		</InLine>
	//	</Ad>
	//</VAST>

	var finalStaticSeats = make([]staticSeat, 0)

	for j := range seats {
		d, err := staticseat.GetStaticSeat(pub, "vast", seats[j].Start)
		if err == nil && d.Chance() > rand.Intn(100) {
			finalStaticSeats = append(finalStaticSeats, staticSeat{
				staticSeat: d,
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

	ua := user_agent.New(r.UserAgent())
	mi := 0
	if ua.Mobile() {
		mi = 1
	}

	rIP := framework.RealIP(r)
	rUserAgent := r.UserAgent()

	bq := &openrtb.BidRequest{
		ID: <-random.ID,
		User: &openrtb.User{
			ID: vastUserIDGenerator(tid, rUserAgent, rIP),
		},
		Imp: imps,
		Site: &openrtb.Site{
			Mobile: mi,
			Page:   l,
			Ref:    ref,
			Inventory: openrtb.Inventory{
				Domain: pub.Name(),
				Name:   pub.Name(),
				ID:     fmt.Sprint(pub.ID()),
				Cat:    pub.Categories(),
			},
		},
		Device: &openrtb.Device{
			IP:  rIP,
			DNT: dnt,
			OS:  ua.OS(),
			UA:  rUserAgent,
		},
	}

	var br = &openrtb.BidResponse{}

	// better since the json is static :)
	bq.Ext = []byte(`{"capping_mode": "reset","underfloor":true}`)

	br, err = client.Call(ctx, method.String(), server.String(), bq)
	if err != nil {
		e := "demand error"
		writesErrorStatus(w, http.StatusInternalServerError, e)
		xlog.GetWithError(ctx, err).Debugf(e)
		return
	}

	for i := range finalStaticSeats {
		seats[finalStaticSeats[i].key] = finalStaticSeats[i].seat[finalStaticSeats[i].key]
		br.SeatBid = append(br.SeatBid, openrtb.SeatBid{
			Bid: []openrtb.Bid{
				{
					ID:       <-random.ID,
					ImpID:    finalStaticSeats[i].key,
					AdMarkup: finalStaticSeats[i].staticSeat.RTBMarkup(),
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

// vastUserIDGenerator create user id for vast
func vastUserIDGenerator(tid, ua, ip string) string {
	return simplehash.MD5(fmt.Sprintf("%s%s%s", tid, ua, ip))
}
