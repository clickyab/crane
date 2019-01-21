package ortb

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/clickyab/services/random"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/filter"
	"clickyab.com/crane/demand/layers/output/demand"
	"clickyab.com/crane/demand/reducer"
	"clickyab.com/crane/demand/rtb"
	"clickyab.com/crane/metrics"
	"clickyab.com/crane/models/apps"
	"clickyab.com/crane/models/suppliers"
	"clickyab.com/crane/models/website"
	openrtb "clickyab.com/crane/openrtb/v2.5"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/safe"
	"github.com/clickyab/services/version"
	"github.com/clickyab/services/xlog"
	"github.com/golang/protobuf/jsonpb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/xmux"
)

const demandPath = "/ortb/:token"

var (
	ortbWebSelector = []reducer.Filter{
		&filter.Desktop{},
		&filter.OS{},
		&filter.WhiteList{},
		&filter.BlackList{},
		&filter.Category{},
		&filter.Province{},
		&filter.ISP{},
		&filter.Strategy{},
		&filter.ReTargeting{},
	}

	ortbAppSelector = []reducer.Filter{
		&filter.AppBrand{},
		&filter.ConnectionType{},
		&filter.AppCarrier{},
		&filter.WhiteList{},
		&filter.BlackList{},
		&filter.Category{},
		&filter.Province{},
		&filter.ISP{},
		&filter.AreaInGlob{},
		&filter.Strategy{},
		&filter.ReTargeting{},
	}
)

func writesErrorStatus(w http.ResponseWriter, status int, detail string) {
	assert.False(status == http.StatusOK)
	w.WriteHeader(status)
	_, _ = fmt.Fprint(w, detail)
}

// openRTBInput is the route for rtb input layer
func openRTBInput(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	tn := time.Now()
	token := xmux.Param(ctx, "token")
	sup, err := suppliers.GetSupplierByToken(token)

	defer func() {
		var supName = "unknown"
		if sup != nil {
			supName = sup.Name()
		}
		metrics.Duration.With(
			prometheus.Labels{
				"sup":   supName,
				"route": "rest",
			},
		).Observe(time.Since(tn).Seconds())

		metrics.CounterRequest.With(prometheus.Labels{
			"sup":   supName,
			"route": "rest",
		}).Inc()
	}()

	if err != nil {
		e := fmt.Sprintf("supplier with token %s not found", token)
		xlog.GetWithError(ctx, err).Debug(e)
		writesErrorStatus(w, http.StatusNotFound, e)
		return
	}
	payload := &openrtb.BidRequest{}

	err = jsonpb.Unmarshal(r.Body, payload)
	defer assert.Nil(r.Body.Close())
	if err != nil {
		xlog.GetWithError(ctx, err).Errorf("invalid request from %s", sup.Name())
		writesErrorStatus(w, http.StatusBadRequest, err.Error())
		return
	}

	var fatFinger,
		prevent,
		underfloor bool
	var capping = openrtb.Capping_Reset
	var strategy []string
	var tiny = sup.TinyMark()

	if payload.Ext != nil {
		fatFinger = payload.Ext.GetFatFinger()
		prevent = payload.Ext.GetPrevent()
		underfloor = payload.Ext.GetUnderfloor()
		capping = payload.Ext.GetCapping()
		strategy = payload.Ext.GetStrategy()
		tiny = payload.Ext.GetTiny()
	}

	if err := validate(payload); err != nil {
		xlog.GetWithError(ctx, err).Errorf("invalid data from %s. payload: %#v", sup.Name(), payload)
		writesErrorStatus(w, http.StatusBadRequest, err.Error())
		return
	}
	var domain, bundle string
	if payload.GetSite() != nil {

		domain = payload.GetSite().GetDomain()
	}
	if payload.GetApp() != nil {
		bundle = payload.GetApp().GetBundle()
	}
	publisher, selector, ps, prevent, err := handlePublisherSelector(domain, bundle, sup, prevent)

	if err != nil {
		e := fmt.Sprintf("publisher from %s,  not supported: %s. payload: %#v", sup.Name(), ps, payload)

		writesErrorStatus(w, http.StatusBadRequest, e)
		xlog.GetWithError(ctx, err).Debug(e)
		return
	}
	proto := entity.HTTP
	for i := range payload.Imp {
		if payload.Imp[i].Secure == 1 {
			proto = entity.HTTPS
			break
		}
	}

	ua := ""
	ip := ""
	if payload.GetDevice() != nil {
		ua = strings.Trim(payload.GetDevice().GetUa(), "\n\t ")
		ip = strings.Trim(payload.GetDevice().GetIp(), "\n\t ")
	}
	if payload.GetUser() == nil {
		payload.User = &openrtb.User{
			Id:   <-random.ID,
			Data: []*openrtb.UserData{},
		}
	}
	us := payload.GetUser().GetId()

	if ua == "" || ip == "" {
		err := fmt.Errorf("no ip/no ua")
		xlog.GetWithError(ctx, err).Debugf("invalid request from %s payload: %#v", sup.Name(), payload)
		writesErrorStatus(w, http.StatusBadRequest, err.Error())
		return
	}
	sh := fmt.Sprintf("CLICK_%x", sha1.Sum([]byte(fmt.Sprintf("%s_%s_%s_%s", prefix, time.Now().Format(format), ip, ua))))
	perHour, _ := strconv.ParseInt(kv.NewEavStore(sh).AllKeys()["C"], 10, 64)
	if perHour > dailyClickLimit.Int64() {
		w.Header().Set("content-type", "application/json")
		j := jsonpb.Marshaler{}
		assert.Nil(j.Marshal(w, &openrtb.BidResponse{
			Id: payload.Id,
		}))
		return
	}

	b := []builder.ShowOptionSetter{
		builder.SetPublisher(publisher),
		builder.SetTimestamp(),
		builder.SetTargetHost(sup.ShowDomain()),
		builder.SetOSUserAgent(ua),
		builder.SetIPLocation(ip, payload.User, payload.Device, sup),
		builder.SetUser(payload.GetUser().GetData()),
		builder.SetProtocol(proto),
		builder.SetTID(us, payload.GetDevice().GetDidsha1()),
		builder.SetNoTiny(!tiny),
		builder.SetBannerMarkup(sup),
		builder.SetFatFinger(fatFinger),
		builder.SetStrategy(strategy, sup),
		builder.SetRate(float64(sup.Rate())),
		builder.SetPreventDefault(prevent),
		builder.SetCappingMode(entity.CappingMode(capping)),
		builder.SetUnderfloor(underfloor),
		builder.SetCategory(payload),
	}

	// TODO : if we need to implement native/app/vast then the next line must be activated and customized
	// b = append(b, builder.SetFloorPercentage(100), builder.SetMinBidPercentage(100))

	b = setPublisherCustomContext(payload, b, publisher)
	sd, vast := seatDetail(payload)
	if vast {
		b = append(b, builder.SetMultiVideo(true))
	}
	b = append(b, builder.SetDemandSeats(sd...))

	c, err := rtb.Select(ctx, selector, b...)
	if err != nil {
		xlog.GetWithError(ctx, err).Errorf("invalid request from %s", sup.Name())

		writesErrorStatus(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := demand.Render(ctx, c, payload.Id)

	defer safe.GoRoutine(ctx, func() {

		for _, s := range sd {
			metrics.Size.With(prometheus.Labels{
				"sup":  sup.Name(),
				"size": s.Size,
				"io":   "in",
			}).Inc()
		}

	})

	w.Header().Set("crane-version", fmt.Sprint(vs.Count))
	w.Header().Set("content-type", "application/json")
	assert.Nil(err)
	j := jsonpb.Marshaler{}
	assert.Nil(j.Marshal(w, res))
}

var vs = version.GetVersion()

func setPublisherCustomContext(payload *openrtb.BidRequest, b []builder.ShowOptionSetter, publisher entity.Publisher) []builder.ShowOptionSetter {
	if payload.GetSite() != nil {
		b = append(b, builder.SetParent(payload.GetSite().GetPage(), payload.GetSite().GetRef()))
	}
	if payload.GetApp() != nil && payload.GetDevice() != nil {
		b = append(b, builder.SetConnType(payload.GetDevice().GetConnectiontype()))
		b = append(b, builder.SetCarrier(payload.GetDevice().GetCarrier(), publisher))
		b = append(b, builder.SetBrand(payload.GetDevice().GetModel()))
	}
	return b
}

func seatDetail(req *openrtb.BidRequest) ([]builder.DemandSeatData, bool) {
	var (
		imp    = req.Imp
		seats  = make([]builder.DemandSeatData, 0)
		w, h   int32
		vast   bool
		assets []*openrtb.NativeRequest_Asset
	)
	for i := range imp {
		var t entity.RequestType
		if imp[i].GetVideo() != nil {
			w, h = imp[i].Video.W, imp[i].Video.H
			t = entity.RequestTypeVast
			// We just support version 3
			ver := false
			for _, pc := range imp[i].Video.Protocols {
				if pc == openrtb.Protocol_VASTX3X0 {
					ver = true
				}
			}
			if !ver {
				continue
			}
		} else if imp[i].GetBanner() != nil {
			w, h = imp[i].GetBanner().W, imp[i].GetBanner().H
			t = entity.RequestTypeBanner
		} else if imp[i].GetNative() != nil {
			t = entity.RequestTypeNative
			assets = imp[i].GetNative().GetRequestNative().GetAssets()
		}

		seats = append(seats, builder.DemandSeatData{
			MinBid: imp[i].GetBidfloor(),
			PubID:  imp[i].Id,
			Size:   fmt.Sprintf("%dx%d", w, h),
			Type:   t,
			Video:  imp[i].GetVideo(),
			Banner: imp[i].Banner,
			Assets: assets,
			MinCPC: func() float64 {
				if ex := imp[i].GetExt(); ex != nil {
					return ex.Mincpc
				}
				return 0
			}(),
		})
	}
	return seats, vast
}

func handlePublisherSelector(domain, bundle string, sup entity.Supplier, prevent bool) (entity.Publisher, []reducer.Filter, string, bool, error) {
	var (
		publisher entity.Publisher
		selector  []reducer.Filter
		ps        string
		err       error
	)
	if domain != "" {
		ps = domain
		publisher, err = website.GetWebSiteOrFake(sup, domain)
		prevent = false // do not accept prevent default on web requestType
		selector = ortbWebSelector

	} else if bundle != "" {
		ps = bundle
		publisher, err = apps.GetAppOrFake(sup, bundle)
		selector = ortbAppSelector

	} else {
		ps = "None"
		err = errors.New("publisher not supported")
	}
	return publisher, selector, ps, prevent, err
}
