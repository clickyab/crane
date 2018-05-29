package web

import (
	"context"
	"encoding/json"
	"net/http"

	"strings"

	"strconv"

	"fmt"

	"errors"

	"math/rand"

	"text/template"

	"bytes"

	"clickyab.com/crane/demand/entity"
	website "clickyab.com/crane/models/clickyabwebsite"
	"clickyab.com/crane/supplier/client"
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
	server = config.RegisterString("crane.supplier.banner.url", "", "route for banner")
	method = config.RegisterString("crane.supplier.banner.method", "POST", "method for banner request")
	showT  = config.RegisterInt64("crane.supplier.showt", 2, "chance of showt")

	templ *template.Template
)

func init() {
	templ = template.Must(template.New("banner").Parse(`<div style="width:{{ .W }}px; height:{{ .H }}px; overflow:hidden; display:inline;" >{{ .Markup }}<iframe src="//t.clickyab.com/" width="1" height="1" frameborder="0"></iframe></div>`))
}

type size struct {
	Width,
	Height int
}

// Sizes contain all allowed size
var sizesModel = map[int]*size{
	1:  {Width: 120, Height: 600},
	2:  {Width: 160, Height: 600},
	3:  {Width: 300, Height: 250},
	4:  {Width: 336, Height: 280},
	5:  {Width: 468, Height: 60},
	6:  {Width: 728, Height: 90},
	7:  {Width: 120, Height: 240},
	8:  {Width: 320, Height: 50},
	9:  {Width: 800, Height: 440},
	11: {Width: 300, Height: 600},
	12: {Width: 970, Height: 90},
	13: {Width: 970, Height: 250},
	14: {Width: 250, Height: 250},
	15: {Width: 300, Height: 1050},
	16: {Width: 320, Height: 480},
	17: {Width: 48, Height: 320},
	18: {Width: 128, Height: 128},
}
var sup = supplier.NewClickyab()

//	d		: domain
//	l		: location
//	r		: ref
//	c		: count of impression. must match with slot count // TODO : do we need it?
//	s		: slots
//	m		: mobile
//	tid		: tracking id
func getAd(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pubID := r.URL.Query().Get("i")
	pub, err := website.GetWebSite(sup, pubID)

	if err != nil {
		xlog.GetWithError(ctx, err).Debugf("website with publisher id %s and supplier %s not found", pubID, sup)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "website not found")
		return
	}
	l := r.URL.Query().Get("l")
	ref := r.URL.Query().Get("r")
	dnt, _ := strconv.Atoi(r.Header.Get("DNT"))
	m, _ := strconv.ParseBool(r.URL.Query().Get("m"))
	tid := r.URL.Query().Get("tid")
	s := r.URL.Query().Get("s")
	c, err := strconv.Atoi(r.URL.Query().Get("c"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid c param")
		return
	}
	_, ok := pub.Attributes()[entity.PAMobileAd]
	extra := ""
	if ok && m {
		extra = simplehash.CRC32(pub.Name())
	}
	imps, err := exSlot(ctx, s, c, r, pub, extra)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

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
			ID: webUserIDGenerator(tid, rUserAgent, rIP),
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

	// better since the json is static :)
	ext := map[string]interface{}{
		"capping_mode": "reset",
		"underfloor":   true,
	}
	// fat finger is allowed only on mobile
	if _, ok := pub.Attributes()[entity.PAFatFinger]; ok && m {
		ext["fat_finger"] = true
	}
	j, err := json.Marshal(ext)
	assert.Nil(err)
	bq.Ext = j
	br, err := client.Call(ctx, method.String(), server.String(), bq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		xlog.GetWithError(ctx, err).Debug("error in call demand server")

		return
	}

	if len(br.SeatBid) > 0 && mi == 1 && rand.Int63n(100) <= showT.Int64() {
		buf := &bytes.Buffer{}
		_ = templ.Execute(buf, struct {
			W, H   int
			Markup string
		}{
			W:      br.SeatBid[0].Bid[0].W,
			H:      br.SeatBid[0].Bid[0].H,
			Markup: br.SeatBid[0].Bid[0].AdMarkup,
		})
		br.SeatBid[0].Bid[0].AdMarkup = buf.String()
	}

	if output.RenderBanner(ctx, w, br, extra) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func exSlot(ctx context.Context, s string, l int, r *http.Request, pub entity.Publisher, extra string) ([]openrtb.Impression, error) {
	sec := secure(r)
	res := make([]openrtb.Impression, 0)
	ts := strings.Split(s, ",")
	if len(ts) != l {
		xlog.Get(ctx).Debug("len of impression does not match with request")
		return nil, errors.New("len of impression does not match with request")
	}
	// calculate min cpc and insert in impression ext
	impExt := map[string]interface{}{
		"min_cpc": pub.MinCPC(string(entity.RequestTypeBanner)),
	}
	iExt, err := json.Marshal(impExt)
	assert.Nil(err)
	for _, v := range ts {
		tv := strings.Split(v, ":")

		if len(tv) != 2 {
			xlog.Get(ctx).Debug("split on s(query param) does not return string array with len of two")

			return nil, errors.New("split on s(query param) does not return string array with len of two")
		}
		sz, err := strconv.Atoi(tv[1])
		if err != nil {
			return nil, err
		}
		var w, h int
		if sizeVal, ok := sizesModel[sz]; ok {
			w, h = sizeVal.Width, sizeVal.Height
		}
		if w == 0 || h == 0 {
			xlog.Get(ctx).Debug("wrong size")
			return nil, errors.New("wrong size")
		}

		res = append(res, openrtb.Impression{
			ID:     tv[0],
			Secure: sec,
			Banner: &openrtb.Banner{
				ID: tv[0],
				H:  h,
				W:  w,
			},
			Ext:      iExt,
			BidFloor: float64(pub.FloorCPM()),
		})

	}
	if extra != "" {
		res = append(res, openrtb.Impression{
			ID:     extra,
			Secure: sec,
			Banner: &openrtb.Banner{
				ID: extra,
				H:  50,
				W:  320,
			},
			BidFloor: float64(pub.FloorCPM()),
			Ext:      iExt,
		})
	}
	return res, nil
}

// webUserIDGenerator create userID for web request
func webUserIDGenerator(tid, ua, ip string) string {
	return simplehash.MD5(fmt.Sprintf("%s%s%s", tid, ua, ip))
}

// secure check openrtb protocol (http/https)
func secure(r *http.Request) openrtb.NumberOrString {
	if framework.Scheme(r) == "https" {
		return openrtb.NumberOrString(1)
	}
	return openrtb.NumberOrString(0)
}
