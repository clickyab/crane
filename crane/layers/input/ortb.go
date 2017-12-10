package input

import (
	"context"
	"encoding/json"
	"net/http"

	"clickyab.com/crane/crane/builder"
	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/filter"
	"clickyab.com/crane/crane/layers/output/demand"
	"clickyab.com/crane/crane/models"
	"clickyab.com/crane/crane/reducer"
	"clickyab.com/crane/crane/rtb"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/xlog"
	"github.com/rs/xmux"
)

var (
	noTiny = config.RegisterBoolean("input.no_tiny", false, "is the boolean case for no tiny")

	webSelector = reducer.Mix(
		filter.IsWebNetwork,
		filter.IsWebMobile,
		filter.CheckDesktopNetwork,
		filter.CheckWebSize,
		filter.CheckOS,
		filter.CheckWhiteList,
		filter.CheckProvince,
		filter.CheckISP,
	)
)

// OrtbInput is the route for rtb input layer
func OrtbInput(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tid := xmux.Param(ctx, "tid")
	payload := openrtb.BidRequest{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("bad request payload"))
		if err != nil {
			xlog.SetField(ctx, "request body", err)
		}
		return
	}

	if !requestValidation(payload) || payload.Validate() != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("not a validate body"))
		if err != nil {
			xlog.SetField(ctx, "request body", err)
		}
		return
	}

	publisher, err := models.GetWebSite(publisher(payload))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, err = w.Write([]byte("website was not found"))
		if err != nil {
			xlog.SetField(ctx, "ioWriter", err)
		}
		return
	}

	contexts := []builder.ShowOptionSetter{
		builder.SetType(requestType(payload)),
		builder.SetOSUserAgent(r.UserAgent()),
		builder.SetIPLocation(framework.RealIP(r)),
		builder.SetPublisher(publisher),
		builder.SetTID(tid),
		builder.SetAlexa(r.UserAgent(), r.Header),
		builder.SetProtocol(r),
		builder.SetQueryParameters(r.URL),
		builder.SetProtocol(r),
		builder.SetNoTiny(noTiny.Bool()),
		builder.SetDemandSeats(seatDetail(payload)),
	}

	conx, err := rtb.Select(ctx, webSelector, contexts...)
	//TODO not a good response or response code, w8ing for merging select
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("supplier name key not found"))
		if err != nil {
			xlog.SetField(ctx, "ioWriter", err)
		}
		return
	}

	header := demand.Rtb(ctx, conx, w)
	for k, v := range header {
		for _, vv := range v {
			w.Header().Set(k, vv)
		}
	}
}

func requestType(req openrtb.BidRequest) entity.RequestType {
	if req.Site != nil {
		if req.Imp[0].Banner != nil {
			// not sure about this one
			return entity.RequestTypeDemand
		} else if req.Imp[0].Native != nil {
			return entity.RequestTypeNative
		}
		return entity.RequestTypeVast
	} else if req.App != nil && req.Imp[0].Banner != nil {
		return entity.RequestTypeApp
	}

	return ""
}

func publisherType(req openrtb.BidRequest) entity.RequestType {
	if req.App != nil {
		return entity.RequestTypeApp
	}
	return entity.RequestTypeWeb
}

func seatDetail(req openrtb.BidRequest) []builder.SeatDetail {
	var (
		imp   = req.Imp
		seats []builder.SeatDetail
		w, h  int
	)

	for i := range imp {
		switch requestType(req) {
		case entity.RequestTypeDemand:
			w, h = imp[i].Banner.W, imp[i].Banner.H
		case entity.RequestTypeApp:
			w, h = imp[i].Banner.W, imp[i].Banner.H
		case entity.RequestTypeVast:
			w, h = imp[i].Video.W, imp[i].Video.H
		}

		seats = append(seats, builder.SeatDetail{
			PubID: imp[i].ID,
			W:     w,
			H:     h,
		})
	}
	return seats
}

func publisher(req openrtb.BidRequest) (name, domain string) {
	if publisherType(req) == entity.RequestTypeWeb {
		name, domain = req.Site.Name, req.Site.Domain
		return
	}
	name, domain = req.App.Name, req.App.Domain
	return
}
