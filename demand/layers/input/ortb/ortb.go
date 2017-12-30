package ortb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"strconv"

	"errors"

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
	ortbSelector = reducer.Mix(
		&filter.WebSize{},
		&filter.WebNetwork{},
		&filter.WebMobile{},
		&filter.Desktop{},
		&filter.OS{},
		&filter.WhiteList{},
		&filter.BlackList{},
		&filter.Category{},
		&filter.Province{},
		&filter.ISP{},
	)
)

type simpleMap map[string]interface{}

func (s simpleMap) Bool(k string) bool {
	d, ok := s[k]
	if !ok {
		return false
	}
	switch t := d.(type) {
	case float64:
		return t != 0
	case string:
		b, _ := strconv.ParseBool(t)
		return b
	case bool:
		return t
	default:
		return false
	}
}

// openrtbInput is the route for rtb input layer
func openrtbInput(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
	var ext = make(simpleMap)
	_ = json.Unmarshal(payload.Ext, &ext)
	fatFinger := ext.Bool("fat_finger")

	if err := payload.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		xlog.GetWithError(ctx, err).Error("invalid data")
		return
	}

	var (
		publisher entity.Publisher
		subType   entity.RequestType
	)
	if payload.Site != nil {
		publisher, err = website.GetWebSite(sup, payload.Site.Domain)
		subType = entity.RequestTypeWeb
	} else if payload.App != nil {
		publisher, err = apps.GetApp(sup, payload.App.Bundle)
		subType = entity.RequestTypeApp
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
		ua = payload.Device.UA
		ip = payload.Device.IP
	}
	us := ""
	if payload.User != nil {
		us = payload.User.ID
	}
	b := []builder.ShowOptionSetter{
		builder.SetTimestamp(),
		builder.SetType(entity.RequestTypeDemand, subType),
		builder.SetTargetHost(sup.ShowDomain()),
		builder.SetOSUserAgent(ua),
		builder.SetIPLocation(ip),
		builder.SetPublisher(publisher),
		builder.SetAlexa(ua, http.Header{}),
		builder.SetProtocol(proto),
		builder.SetTID(us, ip, ua),
		builder.SetNoTiny(sup.TinyMark()),
		builder.SetFatFinger(fatFinger),
		// Website of demand has no floor cpm and soft floor cpm associated with them
		// TODO : decide about this specific values
		builder.SetFloorCPM(sup.DefaultFloorCPM()),
		builder.SetSoftFloorCPM(sup.DefaultSoftFloorCPM()),
		builder.SetRate(float64(sup.Rate())),
	}
	// TODO : if we need to implement native/app/vast then the next line must be activated and customized
	//b = append(b, builder.SetFloorPercentage(100), builder.SetMinBidPercentage(100))

	if payload.Site != nil {
		b = append(b, builder.SetParent(payload.Site.Page, payload.Site.Ref))
	}

	b = append(b, builder.SetDemandSeats(seatDetail(payload)...))

	c, err := rtb.Select(ctx, ortbSelector, b...)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		xlog.GetWithError(ctx, err).Error("invalid request")
		return
	}

	assert.Nil(demand.Render(ctx, w, c))
}

func seatDetail(req openrtb.BidRequest) []builder.DemandSeatData {
	var (
		imp   = req.Imp
		seats = make([]builder.DemandSeatData, 0)
		w, h  int
	)

	for i := range imp {
		if imp[i].Banner != nil {
			w, h = imp[i].Banner.W, imp[i].Banner.H
		}
		seats = append(seats, builder.DemandSeatData{
			MinBid: imp[i].BidFloor,
			PubID:  imp[i].ID,
			Size:   fmt.Sprintf("%dx%d", w, h),
		})
	}
	return seats
}
