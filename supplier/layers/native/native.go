package native

import (
	"context"
	"net/http"
	"strconv"
	"sync"

	"fmt"

	"errors"

	"clickyab.com/crane/models/website"
	"clickyab.com/crane/supplier/client"
	"clickyab.com/crane/supplier/layers/internal/supplier"
	"clickyab.com/crane/supplier/layers/output"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/array"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/simplehash"
	"github.com/clickyab/services/xlog"
	"github.com/mssola/user_agent"
)

var (
	nativeMaxCount        = config.RegisterInt("crane.supplier.native.max_count", 12, "")
	nativeMaxTitleLen     = config.RegisterInt("crane.supplier.native.title,len", 50, "")
	defaultNativeTemplate = config.RegisterString("crane.supplier.native.default.template", "grid4x", "")
	server                = config.RegisterString("crane.supplier.banner.url", "http://127.0.0.1:8090/api/ortb/forbidden", "route for banner")
	method                = config.RegisterString("crane.supplier.banner.method", "POST", "method for banner request")
	template              = make(map[string]temp)
	lock                  = sync.RWMutex{}
)

type temp struct {
	validCounts []int
	image, text bool
}

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
// t			:template (grid3x,grid4x,single,text-list)
// handle supplier native route
func getNative(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("d")
	temp := r.URL.Query().Get("t")
	count := r.URL.Query().Get("count")
	intCount, err := strconv.Atoi(count)
	if intCount > nativeMaxCount.Int() {
		intCount = nativeMaxCount.Int()
	}
	if err != nil || intCount < 1 {
		xlog.GetWithError(ctx, err).Debug("wrong count")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if temp == "" {
		temp = defaultNativeTemplate.String()
	}
	lock.RLock()
	defer lock.RUnlock()
	s, ok := template[temp]
	if !ok {
		xlog.GetWithError(ctx, errors.New("native template is invalid")).Debug("native template is invalid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !array.IntInArray(intCount, s.validCounts...) {
		xlog.GetWithError(ctx, errors.New("wrong native count for template")).Debug("wrong native count for template")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pub, err := website.GetWebSite(sup, domain, 0)
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
		Imp: getImps(r, intCount, pub, s.image, s.text),
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

	if len(br.SeatBid) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	//validate len of response
	br.SeatBid = reduceNativeSeat(br, intCount, s.validCounts)
	result, err := output.RenderNative(ctx, br, !s.image)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = w.Write(result)
	assert.Nil(err)
}

// reduceNativeSeat reduce based on their template configs
func reduceNativeSeat(br *openrtb.BidResponse, intCount int, validCounts []int) []openrtb.SeatBid {
	bidRespCount := len(br.SeatBid)
	if bidRespCount < intCount {
		var valid int
		for i := range validCounts {
			if bidRespCount >= validCounts[i] {
				valid = validCounts[i]
			}
		}
		return br.SeatBid[0:valid]
	}
	return br.SeatBid
}

// registerTemp
func registerTemplate(name string, count []int, image, text bool) {
	lock.Lock()
	defer lock.Unlock()
	_, ok := template[name]
	assert.False(ok)
	template[name] = temp{
		validCounts: count,
		image:       image,
		text:        text,
	}
}

// init register custom clickyab native template
func init() {
	registerTemplate("grid3x", []int{3, 6, 12}, true, true)
	registerTemplate("grid4x", []int{4, 8, 12}, true, true)
	registerTemplate("single", []int{1}, true, true)
	registerTemplate("text-list", []int{1, 3, 4, 6, 8, 12}, false, true)
}

// nativeUserIDGenerator create user id for native
func nativeUserIDGenerator(tid, ua, ip string) string {
	return simplehash.MD5(fmt.Sprintf("%s%s%s", tid, ua, ip))
}
