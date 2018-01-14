package video

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"clickyab.com/crane/models/website"
	"clickyab.com/crane/supplier/client"
	"clickyab.com/crane/supplier/layer/output"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/random"
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
func vast(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	d := r.URL.Query().Get("d")
	pub, err := website.GetWebSite(sup, d)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	l := r.URL.Query().Get("l")
	ref := r.URL.Query().Get("r")
	dnt, _ := strconv.Atoi(r.Header.Get("DNT"))
	m := r.URL.Query().Get("m") != ""
	tid := r.URL.Query().Get("tid")
	ln := r.URL.Query().Get("ln")

	imps, seats := getImps(r, fmt.Sprint(pub.ID()), getSlots(ln))

	ua := user_agent.New(r.UserAgent())
	mi := 0
	if m {
		mi = 1
	}
	bq := &openrtb.BidRequest{
		ID: <-random.ID,
		User: &openrtb.User{
			ID: tid,
		},
		Imp:     imps,
		AllImps: len(imps),
		Site: &openrtb.Site{
			Mobile: mi,
			Page:   l,
			Ref:    ref,
			Inventory: openrtb.Inventory{
				Publisher: &openrtb.Publisher{
					Domain: d,
					Name:   pub.Name(),
					ID:     fmt.Sprint(pub.ID()),
				},
			},
		},
		Device: &openrtb.Device{
			IP:  framework.RealIP(r),
			DNT: dnt,
			OS:  ua.OS(),
			UA:  r.UserAgent(),
		},
	}

	// better since the json is static :)
	bq.Ext = []byte(`{"capping_mode": "reset"}`)
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
