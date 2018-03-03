package video

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"strings"

	website "clickyab.com/crane/models/clickyabwebsite"
	"clickyab.com/crane/supplier/client"
	"clickyab.com/crane/supplier/layer/internal/supplier"
	"clickyab.com/crane/supplier/layer/output"
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

	imps, seats := getImps(r, fmt.Sprint(pub.ID()), getSlots(ln), pub.FloorCPM(), mimes...)

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
		e := "demand error"
		writesErrorStatus(w, http.StatusInternalServerError, e)
		xlog.GetWithError(ctx, err).Debugf(e)
		return
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
