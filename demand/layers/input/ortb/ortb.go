package ortb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"errors"

	"strings"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/filter"
	"clickyab.com/crane/demand/layers/output/demand"
	"clickyab.com/crane/demand/reducer"
	"clickyab.com/crane/demand/rtb"
	apps "clickyab.com/crane/models/apps"
	"clickyab.com/crane/models/suppliers"
	"clickyab.com/crane/models/website"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/xlog"
	"github.com/rs/xmux"
)

const demandPath = "/ortb/:token"

var (
	ortbWebSelector = reducer.Mix(
		&filter.Desktop{},
		&filter.OS{},
		&filter.WhiteList{},
		&filter.BlackList{},
		&filter.Category{},
		&filter.Province{},
		&filter.ISP{},
	)

	ortbAppSelector = reducer.Mix(
		&filter.AppBrand{},
		&filter.AppCarrier{},
		&filter.WhiteList{},
		&filter.BlackList{},
		&filter.Category{},
		&filter.Province{},
		&filter.ISP{},
		&filter.AreaInGlob{},
	)

	ortbVastSelector = reducer.Mix(
		&filter.Desktop{},
		&filter.OS{},
		&filter.WhiteList{},
		&filter.BlackList{},
		&filter.Category{},
		&filter.Province{},
		&filter.ISP{},
	)
)

// openRTBInput is the route for rtb input layer
func openRTBInput(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	token := xmux.Param(ctx, "token")
	sup, err := suppliers.GetSupplierByToken(token)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	payload := openrtb.BidRequest{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&payload); err != nil {
		xlog.GetWithError(ctx, err).Error("invalid request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Known extensions are (currently) fat finger
	var (
		ext         = make(simpleMap)
		cappingMode = entity.CappingStrict
	)
	// If this is not a valid json, just pass by.
	_ = json.Unmarshal(payload.Ext, &ext)
	fatFinger := ext.Bool("fat_finger")
	prevent := ext.Bool("prevent_default")
	capping := ext.String("capping_mode")

	// Currently not supporting no cap (this is intentional)
	if capping == "reset" {
		cappingMode = entity.CappingReset
	}

	if err := payload.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		xlog.GetWithError(ctx, err).Error("invalid data")
		return
	}

	var (
		publisher entity.Publisher
		subType   entity.RequestType
		selector  reducer.Filter
	)
	if payload.Site != nil {
		publisher, err = website.GetWebSite(sup, payload.Site.Domain)
		subType = entity.RequestTypeWeb
		prevent = false // do not accept prevent default on web request
		selector = ortbWebSelector
	} else if payload.App != nil {
		publisher, err = apps.GetApp(sup, payload.App.Bundle)
		subType = entity.RequestTypeApp
		selector = ortbAppSelector
	} else {
		err = errors.New("not supported")
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		xlog.GetWithError(ctx, err).Error("no publisher")
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
		w.WriteHeader(http.StatusNotFound)
		xlog.GetWithError(ctx, err).Error("no ip/no ua")
		return
	}

	b := []builder.ShowOptionSetter{
		builder.SetTimestamp(),
		builder.SetType(entity.RequestTypeDemand, subType),
		builder.SetTargetHost(sup.ShowDomain()),
		builder.SetOSUserAgent(ua),
		builder.SetIPLocation(ip),
		builder.SetPublisher(publisher),
		builder.SetProtocol(proto),
		builder.SetTID(us, ip, ua),
		builder.SetNoTiny(sup.TinyMark()),
		builder.SetFatFinger(fatFinger),

		builder.SetSoftFloorCPM(sup.DefaultSoftFloorCPM()),
		builder.SetRate(float64(sup.Rate())),
		builder.SetPreventDefault(prevent),
		builder.SetCappingMode(cappingMode),
	}
	// TODO : if we need to implement native/app/vast then the next line must be activated and customized
	//b = append(b, builder.SetFloorPercentage(100), builder.SetMinBidPercentage(100))

	if payload.Site != nil {
		b = append(b, builder.SetParent(payload.Site.Page, payload.Site.Ref))
	}
	sd, vast := seatDetail(payload)
	if vast {
		selector = ortbVastSelector
		b = append(b, builder.SetMultiVideo(true))
	}
	b = append(b, builder.SetDemandSeats(sd...))

	c, err := rtb.Select(ctx, selector, b...)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		xlog.GetWithError(ctx, err).Error("invalid request")
		return
	}

	assert.Nil(demand.Render(ctx, w, c))
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
		imp   = req.Imp
		seats = make([]builder.DemandSeatData, 0)
		w, h  int
		vast  bool
	)
	for i := range imp {
		var t builder.SeatType
		if imp[i].Video != nil {
			w, h = imp[i].Video.W, imp[i].Video.H
			t = builder.SeatTypeVideo
			// We just support version 3
			if !intInArray(3, append(imp[i].Video.Protocols, imp[i].Video.Protocol)...) {
				continue
			}
			vast = true
		} else if imp[i].Banner != nil {
			w, h = imp[i].Banner.W, imp[i].Banner.H
			t = builder.SeatTypeBanner
		}
		seats = append(seats, builder.DemandSeatData{
			MinBid: imp[i].BidFloor,
			PubID:  imp[i].ID,
			Size:   fmt.Sprintf("%dx%d", w, h),
			Type:   t,
			Video:  imp[i].Video,
			Banner: imp[i].Banner,
		})
	}
	return seats, vast
}
