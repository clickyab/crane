package native

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"clickyab.com/crane/supplier/middleware/user"

	website "clickyab.com/crane/models/clickyabwebsite"
	openrtb "clickyab.com/crane/openrtb/v2.5"
	"clickyab.com/crane/supplier/client"
	"clickyab.com/crane/supplier/layers/internal/supplier"
	"clickyab.com/crane/supplier/layers/output"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/xlog"
	"github.com/mssola/user_agent"
)

var (
	server            = config.RegisterString("crane.supplier.banner.url", "", "route for banner")
	nativeMaxCount    = config.RegisterInt("crane.supplier.native.max_count", 12, "")
	nativeMaxTitleLen = config.RegisterInt("crane.supplier.native.title,len", 50, "")
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
func getNative(ct context.Context, w http.ResponseWriter, r *http.Request) {
	ctx, cl := context.WithTimeout(ct, time.Second)
	defer cl()
	pubID := r.URL.Query().Get("i")
	pub, err := website.GetWebSite(sup, pubID)
	if err != nil {
		xlog.GetWithError(ctx, err).Debugf("no website: %v", pubID)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("cly-error", fmt.Sprintf("no website: %v", pubID))
		return
	}
	dnt, _ := strconv.Atoi(r.Header.Get("DNT"))
	ref := r.URL.Query().Get("ref")
	parent := r.URL.Query().Get("parent")

	ip := framework.RealIP(r)
	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil || count < 1 {
		xlog.GetWithError(ctx, err).Debug("wrong count")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("cly-error", fmt.Sprintf("wrong count"))

		return
	}
	useragent := r.UserAgent()
	if count > nativeMaxCount.Int() {
		count = nativeMaxCount.Int()
	}

	ua := user_agent.New(useragent)
	var tpl *nativeTemplate

	if ua.Mobile() && count == 3 {
		tpl, err = getNativeTemplate("grid4x")
		count = 2
	} else {
		tpl, err = getNativeTemplate(r.URL.Query().Get("t"))
	}

	if err != nil {
		tpl, err = getNativeTemplate(defaultTemplate.String())
		assert.Nil(err)
	}

	targetCount := getTargetCount(count, tpl.Counts...)
	if targetCount == 0 {
		xlog.GetWithError(ctx, err).Debug("wrong count (during target count calculation)")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("cly-error", fmt.Sprintf("wrong count (during target count calculation)"))

		return
	}
	us, ok := ctx.Value(user.KEY).(*openrtb.User)
	if !ok {
		xlog.GetWithError(ctx, err).Debug("extract user from context")
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("cly-error", "user data")
		return
	}

	fmt.Println(fmt.Sprintf("USER: %#v", us))

	bq := &openrtb.BidRequest{
		Id: fmt.Sprintf("cly-%s", <-random.ID),

		User: us,
		Imp:  getImps(r, targetCount, pub, tpl.Image),
		DistributionchannelOneof: &openrtb.BidRequest_Site{
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

	br, err := client.Call(ctx, server.String(), bq)

	if err != nil {
		xlog.GetWithError(ctx, err).Debugf("Demand: %v ", err)

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
