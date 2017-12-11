package ortb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"clickyab.com/crane/crane/builder"
	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/filter"
	"clickyab.com/crane/crane/layers/output/demand"
	"clickyab.com/crane/crane/models"
	"clickyab.com/crane/crane/reducer"
	"clickyab.com/crane/crane/rtb"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/xlog"
	"github.com/rs/xmux"
)

const demandPath = "/ortb/:token"

var (
	// TODO : verify filters
	webSelector = reducer.Mix(
		filter.IsWebNetwork,
		filter.IsWebMobile,
		filter.CheckDesktopNetwork,
		filter.CheckWebSize,
		filter.CheckOS,
		filter.CheckProvince,
		filter.CheckISP,
	)
)

// openrtbInput is the route for rtb input layer
func openrtbInput(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	token := xmux.Param(ctx, "token")
	sup, err := models.GetSupplierByToken(token)
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

	// TODO : Remove it when app is ready
	// NOTE : This check is not here for so long, so check it again whenever you need to use site
	if payload.Site == nil {
		xlog.Get(ctx).Error("can not support app yet")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := payload.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		xlog.GetWithError(ctx, err).Error("invalid data")
		return
	}

	publisher, err := models.GetWebSite(sup, publisher(payload))
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
	// TODO : verify builders
	b := []builder.ShowOptionSetter{
		builder.SetTimestamp(),
		builder.SetType(entity.RequestTypeDemand),
		builder.SetOSUserAgent(ua),
		builder.SetIPLocation(ip),
		builder.SetPublisher(publisher),
		builder.SetAlexa(ua, http.Header{}),
		builder.SetProtocol(proto),
		builder.SetNoTiny(sup.TinyMark()),
	}

	if payload.Site != nil {
		b = append(b, builder.SetParent(payload.Site.Page, payload.Site.Ref))
	}

	b = append(b, builder.SetDemandSeats(seatDetail(payload)))

	c, err := rtb.Select(ctx, webSelector, b...)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		xlog.GetWithError(ctx, err).Error("invalid request")
		return
	}

	assert.Nil(demand.Render(ctx, w, c))
}

func seatDetail(req openrtb.BidRequest) map[string]string {
	var (
		imp   = req.Imp
		seats = make(map[string]string)
		w, h  int
	)

	for i := range imp {
		if imp[i].Banner != nil {
			w, h = imp[i].Banner.W, imp[i].Banner.H
		}

		seats[imp[i].ID] = fmt.Sprintf("%sx%s", w, h)
	}
	return seats
}

func publisher(req openrtb.BidRequest) string {
	if req.Site != nil {
		return req.Site.Domain
	}
	if req.App != nil {
		return req.App.Domain
	}
	panic("invalid")
}