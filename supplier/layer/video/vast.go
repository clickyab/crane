package video

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"strings"

	website "clickyab.com/crane/models/clickyabwebsite"
	"clickyab.com/crane/supplier/client"
	"clickyab.com/crane/supplier/layer/output"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/simplehash"
	"github.com/mssola/user_agent"
)

var (
	// XXX : currently, there is no need to change the endpoints per type, but if you need it, do it :) its not a rule or something.
	server = config.RegisterString("crane.supplier.banner.url", "", "route for banner")
	method = config.RegisterString("crane.supplier.banner.method", "POST", "method for banner request")
)

var sup = &supplier{}

//	d		: domain
//	l		: location
//	r		: ref
//	ln		: length
//	m		: mobile
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
	m := r.URL.Query().Get("m") != ""
	tid := r.URL.Query().Get("tid")
	ln := r.URL.Query().Get("l")
	var mimes []string
	if mim := strings.Trim(r.URL.Query().Get("mimes"), "\n\t "); mim != "" {
		mimes = strings.Split(mim, ",")
	}

	imps, seats := getImps(r, fmt.Sprint(pub.ID()), getSlots(ln), pub.FloorCPM(), mimes...)

	ua := user_agent.New(r.UserAgent())
	mi := 0
	if m {
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
				Publisher: &openrtb.Publisher{
					Domain: pub.Name(),
					Name:   pub.Name(),
					ID:     fmt.Sprint(pub.ID()),
				},
			},
		},
		Device: &openrtb.Device{
			IP:  rIP,
			DNT: dnt,
			OS:  ua.OS(),
			UA:  rUserAgent,
		},
	}

	// better since the json is static :)
	bq.Ext = []byte(`{"capping_mode": "reset","underfloor":true}`)
	br, err := client.Call(ctx, method.String(), server.String(), bq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if output.RenderVMAP(ctx, w, br, seats) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// vastUserIDGenerator create user id for vast
func vastUserIDGenerator(tid, ua, ip string) string {
	return simplehash.MD5(fmt.Sprintf("%s%s%s", tid, ua, ip))
}
