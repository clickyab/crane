package ortb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/filter"
	"clickyab.com/crane/demand/layers/output/demand"
	"clickyab.com/crane/demand/reducer"
	"clickyab.com/crane/demand/rtb"
	"clickyab.com/crane/models/apps"
	"clickyab.com/crane/models/suppliers"
	"clickyab.com/crane/models/website"
	"github.com/bsm/openrtb"
	"github.com/bsm/openrtb/native/request"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/xlog"
	"github.com/rs/xmux"
)

const demandPath = "/ortb/:token"

var (
	ortbWebSelector = reducer.Mix(
		&filter.Strategy{},
		&filter.Desktop{},
		&filter.OS{},
		&filter.WhiteList{},
		&filter.BlackList{},
		&filter.Category{},
		&filter.Province{},
		&filter.ISP{},
	)

	ortbAppSelector = reducer.Mix(
		&filter.Strategy{},
		&filter.AppBrand{},
		&filter.ConnectionType{},
		&filter.AppCarrier{},
		&filter.WhiteList{},
		&filter.BlackList{},
		&filter.Category{},
		&filter.Province{},
		&filter.ISP{},
		&filter.AreaInGlob{},
	)
)

func writesErrorStatus(w http.ResponseWriter, status int, detail string) {
	assert.False(status == http.StatusOK)
	w.WriteHeader(status)
	fmt.Fprint(w, detail)
}

// openRTBInput is the route for rtb input layer
func openRTBInput(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	token := xmux.Param(ctx, "token")
	sup, err := suppliers.GetSupplierByToken(token)
	if err != nil {
		e := fmt.Sprintf("supplier with token %s not found", token)
		xlog.GetWithError(ctx, err).Debug(e)
		writesErrorStatus(w, http.StatusNotFound, e)
		return
	}
	payload := openrtb.BidRequest{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&payload); err != nil {
		xlog.GetWithError(ctx, err).Error("invalid request")
		writesErrorStatus(w, http.StatusBadRequest, err.Error())
		return
	}

	// Known extensions are (currently) fat finger
	var (
		ext         = make(simpleMap)
		cappingMode = entity.CappingStrict
	)
	// If this is not a valid json, just pass by.
	_ = json.Unmarshal(payload.Ext, &ext)
	fatFinger, _ := ext.Bool("fat_finger")
	prevent, _ := ext.Bool("prevent_default")
	underfloor, _ := ext.Bool("underfloor")
	capping, _ := ext.String("capping_mode")
	var strategy []string
	sts, _ := ext.String("strategy")
	if st := strings.Trim(sts, "\t \n"); st != "" {
		strategy = strings.Split(st, ",")
	}
	tiny, ok := ext.Bool("tiny_mark")
	if !ok {
		tiny = sup.TinyMark()
	}

	// Currently not supporting no cap (this is intentional)
	if capping == "reset" {
		cappingMode = entity.CappingReset
	}

	if err := payload.Validate(); err != nil {
		xlog.GetWithError(ctx, err).Error("invalid data")
		writesErrorStatus(w, http.StatusBadRequest, err.Error())
		return
	}

	publisher, selector, ps, prevent, err := handlePublisherSelector(payload, sup, prevent)

	if err != nil {
		e := fmt.Sprintf("publisher not supported : %s", ps)
		writesErrorStatus(w, http.StatusBadRequest, e)
		xlog.GetWithError(ctx, err).Debug(e)
		return
	}
	proto := entity.HTTP
	for i := range payload.Imp {
		if payload.Imp[i].Secure != 0 {
			proto = entity.HTTPS
			break
		}
	}

	ua := ""
	ip := ""
	if payload.Device != nil {
		ua = strings.Trim(payload.Device.UA, "\n\t ")
		ip = strings.Trim(payload.Device.IP, "\n\t ")
	}
	us := ""
	if payload.User != nil {
		us = payload.User.ID
	}

	if ua == "" || ip == "" {
		err := fmt.Errorf("no ip/no ua")
		xlog.GetWithError(ctx, err).Debug("invalid request")
		writesErrorStatus(w, http.StatusBadRequest, err.Error())
		return
	}

	b := []builder.ShowOptionSetter{
		builder.SetTimestamp(),
		builder.SetTargetHost(sup.ShowDomain()),
		builder.SetOSUserAgent(ua),
		builder.SetIPLocation(ip, payload.User, payload.Device),
		builder.SetPublisher(publisher),
		builder.SetProtocol(proto),
		builder.SetTID(us, ip, ua),
		builder.SetNoTiny(!tiny),
		builder.SetBannerMarkup(sup),
		builder.SetFatFinger(fatFinger),
		builder.SetStrategy(strategy, sup),
		builder.SetRate(float64(sup.Rate())),
		builder.SetPreventDefault(prevent),
		builder.SetCappingMode(cappingMode),
		builder.SetUnderfloor(underfloor),
		builder.SetCategory(&payload),
	}
	// TODO : if we need to implement native/app/vast then the next line must be activated and customized
	//b = append(b, builder.SetFloorPercentage(100), builder.SetMinBidPercentage(100))

	b = setPublisherCustomContext(payload, b)
	sd, vast := seatDetail(payload)
	if vast {
		b = append(b, builder.SetMultiVideo(true))
	}
	b = append(b, builder.SetDemandSeats(sd...))

	c, err := rtb.Select(ctx, selector, b...)
	if err != nil {
		xlog.GetWithError(ctx, err).Error("invalid request")
		writesErrorStatus(w, http.StatusBadRequest, err.Error())
		return
	}

	assert.Nil(demand.Render(ctx, w, c))
}

func setPublisherCustomContext(payload openrtb.BidRequest, b []builder.ShowOptionSetter) []builder.ShowOptionSetter {
	if payload.Site != nil {
		b = append(b, builder.SetParent(payload.Site.Page, payload.Site.Ref))
	}
	if payload.App != nil {
		b = append(b, builder.SetConnType(payload.Device.ConnType))
		b = append(b, builder.SetCarrier(payload.Device.Carrier))
		b = append(b, builder.SetBrand(payload.Device.Model))
	}
	return b
}

func intInArray(v int, all ...int) bool {
	for i := range all {
		if v == all[i] {
			return true
		}
	}

	return false
}

func seatDetail(req openrtb.BidRequest) ([]builder.DemandSeatData, bool) {
	var (
		imp    = req.Imp
		seats  = make([]builder.DemandSeatData, 0)
		w, h   int
		vast   bool
		assets []request.Asset
	)
	for i := range imp {
		var t entity.RequestType
		if imp[i].Video != nil {
			w, h = imp[i].Video.W, imp[i].Video.H
			t = entity.RequestTypeVast
			// We just support version 3
			if !intInArray(3, append(imp[i].Video.Protocols, imp[i].Video.Protocol)...) {
				continue
			}
			vast = true
		} else if imp[i].Banner != nil {
			w, h = imp[i].Banner.W, imp[i].Banner.H
			t = entity.RequestTypeBanner
		} else if imp[i].Native != nil {
			t = entity.RequestTypeNative
			req := request.Request{}
			err := json.Unmarshal(imp[i].Native.Request, &req)
			assert.Nil(err)
			assets = req.Assets
		}
		var (
			ext = make(simpleMap)
		)
		// If this is not a valid json, just pass by.
		_ = json.Unmarshal(imp[i].Ext, &ext)
		seats = append(seats, builder.DemandSeatData{
			MinBid: imp[i].BidFloor,
			PubID:  imp[i].ID,
			Size:   fmt.Sprintf("%dx%d", w, h),
			Type:   t,
			Video:  imp[i].Video,
			Banner: imp[i].Banner,
			Assets: assets,
			MinCPC: ext.Float64("min_cpc"),
		})
	}
	return seats, vast
}

func handlePublisherSelector(payload openrtb.BidRequest, sup entity.Supplier, prevent bool) (entity.Publisher, reducer.Filter, string, bool, error) {
	var (
		publisher entity.Publisher
		selector  reducer.Filter
		ps        string
		err       error
	)
	if payload.Site != nil {
		if payload.Site.Domain == "" {
			err = errors.New("website domain is empty")
		} else {
			ps = payload.Site.Domain
			publisher, err = website.GetWebSiteOrFake(sup, payload.Site.Domain)
			prevent = false // do not accept prevent default on web requestType
			selector = ortbWebSelector
		}
	} else if payload.App != nil {
		if payload.App.Bundle == "" {
			err = errors.New("app bundle is empty")
		} else {
			ps = payload.App.Bundle
			publisher, err = apps.GetAppOrFake(sup, payload.App.Bundle)
			selector = ortbAppSelector
		}

	} else {
		ps = "None"
		err = errors.New("not supported")
	}
	return publisher, selector, ps, prevent, err
}
