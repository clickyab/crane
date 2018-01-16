package web

import (
	"context"
	"net/http"

	"strings"

	"strconv"

	"fmt"

	"hash/crc32"

	"errors"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/models/website"
	"clickyab.com/crane/supplier/client"
	"clickyab.com/crane/supplier/layer/output"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/xlog"
	"github.com/mssola/user_agent"
)

var (
	server = config.RegisterString("crane.supplier.banner.url", "", "route for banner")
	method = config.RegisterString("crane.supplier.banner.method", "POST", "method for banner request")
)

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
var sup entity.Supplier = &supplier{}

//	d		: domain
//	l		: location
//	r		: ref
//	c		: count of impression. must match with slot count // TODO : do we need it?
//	s		: slots
//	m		: mobile
//	tid		: tracking id
func getAd(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
	s := r.URL.Query().Get("s")
	c, err := strconv.Atoi(r.URL.Query().Get("c"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, ok := pub.Attributes()[entity.PAMobileAd]
	extra := ""
	if ok && m {
		extra = crc(d)
	}
	imps, err := exSlot(ctx, s, c, r, extra)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

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
		Imp: imps,
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

	if output.RenderBanner(ctx, w, br, extra) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func crc(d string) string {
	x := crc32.New(crc32.IEEETable)
	_, _ = x.Write([]byte(d))
	return fmt.Sprint(x.Sum32())
}

func exSlot(ctx context.Context, s string, l int, r *http.Request, extra string) ([]openrtb.Impression, error) {
	sec := secure(r)
	res := make([]openrtb.Impression, 0)
	ts := strings.Split(s, ",")
	if len(ts) != l {
		xlog.Get(ctx).Debug("len of impression does not match with request")
		return nil, errors.New("len of impression does not match with request")
	}
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
		w, h := sizesModel[sz].Width, sizesModel[sz].Height
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
		})
	}
	return res, nil
}

func secure(r *http.Request) openrtb.NumberOrString {
	if framework.Scheme(r) == "https" {
		return openrtb.NumberOrString(1)
	}
	return openrtb.NumberOrString(0)
}
