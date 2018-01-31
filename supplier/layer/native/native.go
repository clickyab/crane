package native

import (
	"context"
	"net/http"
	"strconv"

	"fmt"

	"encoding/json"

	"clickyab.com/crane/models/website"
	"clickyab.com/crane/supplier/client"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/simplehash"
	"github.com/mssola/user_agent"
	"github.com/sirupsen/logrus"
)

var (
	nativeMaxCount    = config.RegisterInt("crane.supplier.native.max_count", 12, "")
	nativeMaxTitleLen = config.RegisterInt("crane.supplier.native.title,len", 50, "")
	server            = config.RegisterString("crane.supplier.native.url", "", "route for banner")
	method            = config.RegisterString("crane.supplier.native.method", "POST", "method for banner request")
)

// ImageType openrtb native image
type ImageType int

const (
	// IconImageType icon image
	IconImageType ImageType = 1
	// MainImageType main image
	MainImageType ImageType = 3
)

// d			:domain
// ref			:referrer
// parent		:parent
// count		:number of impression
// m			:mobile or not
// handle supplier native route
func getNative(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("d")
	pub, err := website.GetWebSiteOrFake(&supplier{}, domain)
	if err != nil {
		// TODO send proper message
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dnt, _ := strconv.Atoi(r.Header.Get("DNT"))
	ref := r.URL.Query().Get("ref")
	parent := r.URL.Query().Get("parent")
	tid := r.URL.Query().Get("tid")
	m := r.URL.Query().Get("m") != ""

	ip := framework.RealIP(r)
	useragent := r.UserAgent()

	count := r.URL.Query().Get("count")
	intCount, err := strconv.Atoi(count)
	if err != nil || intCount < 1 {
		// TODO send proper message
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if intCount > nativeMaxCount.Int() {
		intCount = nativeMaxCount.Int()
	}

	ua := user_agent.New(useragent)

	mi := 0
	if m {
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
				Publisher: &openrtb.Publisher{
					Domain: domain,
					Name:   pub.Name(),
					ID:     fmt.Sprint(pub.ID()),
				},
			},
		},
		Device: &openrtb.Device{
			IP:  ip,
			OS:  ua.OS(),
			UA:  useragent,
			DNT: dnt,
		},
	}

	jj, _ := json.Marshal(bq)
	logrus.Debug(string(jj))

	bq.Ext = []byte(`{"capping_mode": "reset"}`)
	br, err := client.Call(ctx, method.String(), server.String(), bq)
	if err != nil {
		// TODO send proper message
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logrus.Debug(br)

}

// nativeUserIDGenerator create user id for native
func nativeUserIDGenerator(tid, ua, ip string) string {
	return simplehash.MD5(fmt.Sprintf("%s%s%s", tid, ua, ip))
}
