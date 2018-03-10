package native

import (
	"context"
	"net/http"
	"strconv"

	"fmt"

	website "clickyab.com/crane/models/clickyabwebsite"
	"clickyab.com/crane/supplier/client"
	"clickyab.com/crane/supplier/layers/internal/supplier"
	"clickyab.com/crane/supplier/layers/output"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/simplehash"
	"github.com/clickyab/services/xlog"
	"github.com/mssola/user_agent"
)

var (
	nativeMaxCount    = config.RegisterInt("crane.supplier.native.max_count", 12, "")
	nativeMaxTitleLen = config.RegisterInt("crane.supplier.native.title,len", 50, "")
	server            = config.RegisterString("crane.supplier.banner.url", "", "route for banner")
	method            = config.RegisterString("crane.supplier.banner.method", "POST", "method for banner request")
)

// ImageType openrtb native image
type ImageType int

const (
	// IconImageType icon image
	IconImageType ImageType = 1
	// MainImageType main image
	MainImageType ImageType = 3
)

var (
	sup = supplier.NewClickyab()
)

// d			:domain
// ref			:referrer
// parent		:parent
// count		:number of impression
// handle supplier native route
func getNative(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pubID := r.URL.Query().Get("a")
	pub, err := website.GetWebSite(sup, pubID)
	if err != nil {
		xlog.GetWithError(ctx, err).Debug("no website")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dnt, _ := strconv.Atoi(r.Header.Get("DNT"))
	ref := r.URL.Query().Get("ref")
	parent := r.URL.Query().Get("parent")
	tid := r.URL.Query().Get("tid")

	ip := framework.RealIP(r)
	useragent := r.UserAgent()

	count := r.URL.Query().Get("count")
	intCount, err := strconv.Atoi(count)
	if err != nil || intCount < 1 {
		xlog.GetWithError(ctx, err).Debug("wrong count")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if intCount > nativeMaxCount.Int() {
		intCount = nativeMaxCount.Int()
	}

	ua := user_agent.New(useragent)

	mi := 0
	if ua.Mobile() {
		mi = 1
	}

	bq := &openrtb.BidRequest{
		ID: <-random.ID,
		User: &openrtb.User{
			ID: nativeUserIDGenerator(tid, useragent, ip),
		},
		Imp: getImps(r, intCount, pub.ID(), pub.FloorCPM()),
		Site: &openrtb.Site{
			Page:   parent,
			Ref:    ref,
			Mobile: mi,
			Inventory: openrtb.Inventory{
				Domain: pub.Name(),
				Name:   pub.Name(),
				ID:     fmt.Sprint(pub.ID()),
			},
		},
		Device: &openrtb.Device{
			IP:  ip,
			OS:  ua.OS(),
			UA:  useragent,
			DNT: dnt,
		},
	}

	bq.Ext = []byte(`{"capping_mode": "reset", "underfloor": true}`)
	br, err := client.Call(ctx, method.String(), server.String(), bq)
	if err != nil {
		// TODO send proper message
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if output.RenderNative(ctx, w, br) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// nativeUserIDGenerator create user id for native
func nativeUserIDGenerator(tid, ua, ip string) string {
	return simplehash.MD5(fmt.Sprintf("%s%s%s", tid, ua, ip))
}
