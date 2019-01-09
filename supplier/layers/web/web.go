package web

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"clickyab.com/crane/demand/entity"
	website "clickyab.com/crane/models/clickyabwebsite"
	openrtb "clickyab.com/crane/openrtb/v2.5"
	"clickyab.com/crane/supplier/client"
	"clickyab.com/crane/supplier/layers/internal/supplier"
	"clickyab.com/crane/supplier/layers/output"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/simplehash"
	"github.com/clickyab/services/xlog"
	"github.com/mssola/user_agent"
)

var (
	showT = config.RegisterInt64("crane.supplier.showt", 100, "chance of showt")

	templ *template.Template
)

func init() {
	templ = template.Must(template.New("banner").Parse(`<div class="cyb-bholder" style="position:relative;width:{{ .W }}px; height:{{ .H }}px; overflow:hidden;" ><style>
.tiny {height: 18px;width: 19px;position: absolute;top: 0px;right: 0;z-index: 100;background: url(//static.clickyab.com/img/clickyab-tiny.png) right top no-repeat;border-bottom-left-radius: 4px;-moz-border-radius-bottomleft: 4px;} .cyb-bholder img {width:100%;}
</style><a class="tiny" href="https://www.clickyab.com" target="_blank"></a>{{ .Markup }} <iframe src="//t.clickyab.com" frameborder="0" height="1" width="1" style="position: absolute; top: -10000px; left: -10000px"></iframe></div>`))

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
	21: {Width: 1200, Height: 627},
}
var sup = supplier.NewClickyab()

var server = config.RegisterString("crane.supplier.banner.url", "", "route for banner")

// 	d		: domain
// 	l		: location
// 	r		: ref
// 	c		: count of impression. must match with slot count // TODO : do we need it?
// 	s		: slots
// 	m		: mobile
// 	tid		: tracking id
func getAd(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pubID := r.URL.Query().Get("i")
	pub, err := website.GetWebSite(sup, pubID)

	if err != nil {
		xlog.GetWithError(ctx, err).Debugf("website with publisher id %s and supplier %s not found", pubID, sup)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("cly-error", fmt.Sprintf("render failed: %s, %s ", pubID, sup.Name()))

		_, _ = fmt.Fprint(w, "website not found")
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
		_, _ = fmt.Fprint(w, "invalid c param")
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
		Id: fmt.Sprintf("cyl-%s", <-random.ID),
		User: &openrtb.User{
			Id: webUserIDGenerator(tid, rUserAgent, rIP),
		},
		Imp: imps,
		DistributionchannelOneof: &openrtb.BidRequest_Site{
			Site: &openrtb.Site{
				Mobile: func() int32 {
					if m {
						return 1
					}
					return 0
				}(),
				Page:   l,
				Ref:    ref,
				Domain: pub.Name(),
				Name:   pub.Name(),
				Id:     fmt.Sprint(pub.ID()),
				Cat:    pub.Categories(),
			},
		},
		Device: &openrtb.Device{
			Ip:  rIP,
			Dnt: int32(dnt),
			Os:  ua.OS(),
			Ua:  rUserAgent,
		},
		Ext: &openrtb.BidRequest_Ext{
			Underfloor: true,
			Capping:    openrtb.Capping_Reset,
			FatFinger: func() bool {
				if _, ok := pub.Attributes()[entity.PAFatFinger]; ok && m {
					return true
				}
				return false
			}(),
		},
	}

	br, err := client.Call(ctx, server.String(), bq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		xlog.GetWithError(ctx, err).Debug("error in call demand server")

		return
	}

	for i := 0; i < len(br.GetSeatbid()); i++ {
		buf := &bytes.Buffer{}
		_ = templ.Execute(buf, struct {
			ShowT  bool
			W, H   int32
			Markup string
		}{
			ShowT:  mi == 1 && rand.Int63n(100) <= showT.Int64() && i == 0,
			W:      br.GetSeatbid()[i].GetBid()[0].GetW(),
			H:      br.GetSeatbid()[i].GetBid()[0].GetH(),
			Markup: br.GetSeatbid()[i].GetBid()[0].GetAdm(),
		})
		br.GetSeatbid()[i].GetBid()[i].AdmOneof = &openrtb.BidResponse_SeatBid_Bid_Adm{
			Adm: buf.String(),
		}
	}

	if output.RenderBanner(ctx, w, br, extra) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func exSlot(ctx context.Context, s string, l int, r *http.Request, pub entity.Publisher, extra string) ([]*openrtb.Imp, error) {
	sec := framework.Scheme(r) == "https"
	res := make([]*openrtb.Imp, 0)
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
		var w, h int
		if sizeVal, ok := sizesModel[sz]; ok {
			w, h = sizeVal.Width, sizeVal.Height
		}
		if w == 0 || h == 0 {
			xlog.Get(ctx).Debug("wrong size")
			return nil, errors.New("wrong size")
		}

		res = append(res, &openrtb.Imp{
			Id: tv[0],
			Secure: func() int32 {
				if sec {
					return 1
				}
				return 0
			}(),
			Banner: &openrtb.Banner{
				Id: tv[0],
				H:  int32(h),
				W:  int32(w),
			},
			Ext: &openrtb.Imp_Ext{
				Mincpc: pub.MinCPC(string(entity.RequestTypeBanner)),
			},
			Bidfloor: float64(pub.FloorCPM()),
		})

	}
	if extra != "" {
		res = append(res, &openrtb.Imp{
			Id: extra,
			Secure: func() int32 {
				if sec {
					return 1
				}
				return 0
			}(),
			Banner: &openrtb.Banner{
				Id: extra,
				H:  50,
				W:  320,
			},
			Bidfloor: float64(pub.FloorCPM()),
			Ext: &openrtb.Imp_Ext{
				Mincpc: pub.MinCPC(string(entity.RequestTypeBanner)),
			},
		})
	}
	return res, nil
}

// webUserIDGenerator create userID for web request
func webUserIDGenerator(tid, ua, ip string) string {
	return simplehash.MD5(fmt.Sprintf("%s%s%s", tid, ua, ip))
}
