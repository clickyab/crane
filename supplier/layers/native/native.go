package native

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	website "clickyab.com/crane/models/clickyabwebsite"
	"clickyab.com/crane/openrtb"
	"clickyab.com/crane/supplier/client"
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
	nativeMaxCount    = config.RegisterInt("crane.supplier.native.max_count", 12, "")
	nativeMaxTitleLen = config.RegisterInt("crane.supplier.native.title,len", 50, "")
	server            = config.RegisterString("crane.supplier.banner.url", "", "route for banner")
	method            = config.RegisterString("crane.supplier.banner.method", "POST", "method for banner request")
	defaultTemplate   = config.RegisterString("crane.supplier.native.default.template", "grid4x", "")
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
// t			:template (grid3x,grid4x,single,text-list)
// ref			:referrer
// parent		:parent
// count		:number of impression
// handle supplier native route
func getNative(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pubID := r.URL.Query().Get("i")
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
	tpl, err := getNativeTemplate(r.URL.Query().Get("t"))
	if err != nil {
		tpl, err = getNativeTemplate(defaultTemplate.String())
		assert.Nil(err)
	}

	ip := framework.RealIP(r)
	useragent := r.UserAgent()

	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil || count < 1 {
		xlog.GetWithError(ctx, err).Debug("wrong count")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if count > nativeMaxCount.Int() {
		count = nativeMaxCount.Int()
	}

	targetCount := getTargetCount(count, tpl.Counts...)
	if targetCount == 0 {
		xlog.GetWithError(ctx, err).Debug("wrong count (during target count calculation)")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ua := user_agent.New(useragent)

	bq := &openrtb.BidRequest{
		Id: fmt.Sprintf("cly-%s", <-random.ID),
		User: &openrtb.User{
			Id: nativeUserIDGenerator(tid, useragent, ip),
		},
		Imp: getImps(r, targetCount, pub, tpl.Image),
		DistributionchannelXoneof: &openrtb.BidRequest_Site{
			Site: &openrtb.Site{
				Page: parent,
				Ref:  ref,
				Mobile: func() int32 {
					if ua.Mobile() {
						return 1
					}
					return 0
				}(),
				Domain: pub.Name(),
				Name:   pub.Name(),
				Id:     fmt.Sprint(pub.ID()),
				Cat:    pub.Categories(),
			},
		},
		Device: &openrtb.Device{
			Ip:  ip,
			Os:  ua.OS(),
			Ua:  useragent,
			Dnt: int32(dnt),
		},
		Ext: &openrtb.BidRequest_Ext{
			Capping:    openrtb.Capping_Reset,
			Underfloor: true,
		},
	}

	br, err := client.Call(ctx, method.String(), server.String(), bq)
	if err != nil {
		// TODO send proper message
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	targetCount = getTargetCount(len(br.GetSeatbid()), tpl.Counts...)

	if targetCount == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	br.Seatbid = br.Seatbid[:targetCount] // drop unwanted count
	result, err := output.RenderNative(ctx, br, tpl.Template)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = w.Write(result)
	assert.Nil(err)
}

func getTargetCount(max int, counts ...int) int {
	target := 0
	for i := range counts {
		if max < counts[i] {
			break
		}
		target = counts[i]
	}

	return target
}

// nativeUserIDGenerator create user id for native
func nativeUserIDGenerator(tid, ua, ip string) string {
	return simplehash.MD5(fmt.Sprintf("%s%s%s", tid, ua, ip))
}
