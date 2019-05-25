package ortb

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/clickyab/services/random"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/layers/output/demand"
	"clickyab.com/crane/demand/rtb"
	"clickyab.com/crane/metrics"
	"clickyab.com/crane/models/suppliers"
	openrtb "clickyab.com/crane/openrtb/v2.5"
	"github.com/clickyab/services/safe"
	"github.com/clickyab/services/xlog"
	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Validation errors
var (
	ErrInvalidReqNoID        = errors.New("openrtb: request ID missing")
	ErrInvalidReqNoImps      = errors.New("openrtb: request has no impressions")
	ErrInvalidReqMultiInv    = errors.New("openrtb: request has no inventory sources") // has site and app
	ErrInvalidImpNoID        = errors.New("openrtb: impression ID missing")
	ErrInvalidImpMultiAssets = errors.New("openrtb: impression has multiple assets") // at least two out of Banner, Video, Native

	ErrInvalidVideoNoMimes       = errors.New("openrtb: video has no mimes")
	ErrInvalidVideoNoLinearity   = errors.New("openrtb: video linearity missing")
	ErrInvalidVideoNoMinDuration = errors.New("openrtb: video min-duration missing")
	ErrInvalidVideoNoMaxDuration = errors.New("openrtb: video max-duration missing")
	ErrInvalidVideoNoProtocols   = errors.New("openrtb: video protocols missing")
)

func assetCount(imp *openrtb.Imp) int {
	n := 0
	if imp.GetBanner() != nil {
		n++
	}
	if imp.GetVideo() != nil {
		n++
	}
	if imp.GetNative() != nil {
		n++
	}
	return n
}

func videoValidate(v *openrtb.Video) error {
	if len(v.GetMimes()) == 0 {
		return ErrInvalidVideoNoMimes
	} else if v.GetLinearity() == 0 {
		return ErrInvalidVideoNoLinearity
	} else if v.GetMinduration() == 0 {
		return ErrInvalidVideoNoMinDuration
	} else if v.GetMaxduration() == 0 {
		return ErrInvalidVideoNoMaxDuration
	} else if len(v.GetProtocols()) == 0 {
		return ErrInvalidVideoNoProtocols
	}
	return nil
}

func validate(req *openrtb.BidRequest) error {
	if req.GetId() == "" {
		return ErrInvalidReqNoID
	} else if req.Imp == nil {
		return ErrInvalidReqNoImps

	}

	for _, imp := range req.GetImp() {
		if imp.GetId() == "" {
			return ErrInvalidImpNoID
		}
		if assetCount(imp) > 1 {
			return ErrInvalidImpMultiAssets
		}
		if imp.GetVideo() != nil {
			if err := videoValidate(imp.GetVideo()); err != nil {
				return err
			}
		}
	}

	return nil
}

// GrpcHandler for handling openrtb request
func GrpcHandler(ctx context.Context, req *openrtb.BidRequest) (*openrtb.BidResponse, error) {
	tn := time.Now()
	xlog.GetWithField(ctx, "REQUEST FROM SUPPLIER", req.Id).Debug()

	res := &openrtb.BidResponse{}
	res.Id = req.GetId()
	token := req.GetToken()
	sup, err := suppliers.GetSupplierByToken(token)
	if err != nil {
		e := fmt.Errorf("supplier with token %s not found", token)
		xlog.GetWithError(ctx, err).Debug(e)
		return nil, e
	}
	defer func() {
		var supName = "unknown"
		if sup != nil {
			supName = sup.Name()
		}
		metrics.Duration.With(
			prometheus.Labels{
				"sup":   supName,
				"route": "grpc",
			},
		).Observe(time.Since(tn).Seconds())

		metrics.CounterRequest.With(prometheus.Labels{
			"sup":   supName,
			"route": "grpc",
		}).Inc()
	}()

	var (
		tiny, fatFinger, prevent, underfloor bool
		strategy                             []string
		capping                              = entity.CappingReset
	)
	if req.Ext != nil {
		req.GetExt()
		fatFinger = req.Ext.GetFatFinger()
		prevent = req.Ext.GetPrevent()
		underfloor = req.Ext.GetUnderfloor()
		tiny = req.Ext.GetTiny()
		strategy = req.Ext.GetStrategy()
	}

	if err := validate(req); err != nil {
		return res, err
	}

	var domain, bundle string
	if req.GetSite() != nil {
		domain = req.GetSite().GetDomain()
	}
	if req.GetApp() != nil {
		bundle = req.GetApp().GetBundle()
	}

	pub, selector, ps, prevent, err := handlePublisherSelector(domain, bundle, sup, prevent)

	if err != nil {
		fmt.Println(err)
		e := spew.Sprintf("publisher from %s, %s, %s, not supported: %s. payload: %#v", sup.Name(), ps, req)
		xlog.GetWithError(ctx, err).Debug(e)
		return nil, grpc.ErrServerStopped
	}

	proto := entity.HTTP
	for i := range req.Imp {
		if req.Imp[i].Secure == 1 {
			proto = entity.HTTPS
			break
		}
	}

	ua := ""
	ip := ""
	if req.GetDevice() != nil {
		ua = strings.Trim(req.GetDevice().GetUa(), "\n\t ")
		ip = strings.Trim(req.GetDevice().GetIp(), "\n\t ")
	}

	if req.GetUser() == nil {
		req.User = &openrtb.User{
			Id:   <-random.ID,
			Data: []*openrtb.UserData{},
		}
	}
	us := req.GetUser().GetId() + req.GetUser().GetBuyeruid()

	if ua == "" || ip == "" {
		err := fmt.Errorf("no ip/no ua")
		xlog.GetWithError(ctx, err).Debug("invalid request")
		return res, err
	}
	b := []builder.ShowOptionSetter{
		builder.SetTimestamp(),
		builder.SetTargetHost(sup.ShowDomain()),
		builder.SetOSUserAgent(ua),
		builder.SetIPLocation(ip, req.GetUser(), req.GetDevice(), sup),
		builder.SetPublisher(pub),
		builder.SetProtocol(proto),
		builder.SetTID(us, req.GetDevice().GetDidsha1()),
		builder.SetUser(req.GetUser().GetData()),
		builder.SetNoTiny(!tiny),
		builder.SetBannerMarkup(sup),
		builder.SetFatFinger(fatFinger),
		builder.SetStrategy(strategy, sup),
		builder.SetRate(float64(sup.Rate())),
		builder.SetPreventDefault(prevent),
		builder.SetCappingMode(capping),
		builder.SetUnderfloor(underfloor),
		builder.SetCategory(req),
	}
	// TODO : if we need to implement native/app/vast then the next line must be activated and customized
	// b = append(b, builder.SetFloorPercentage(100), builder.SetMinBidPercentage(100))

	b = setPublisherCustomContext(req, b, pub)
	sd, vast, version := seatDetail(req)
	if vast {
		b = append(b, builder.SetMultiVideo(true))
	}
	b = append(b, builder.SetDemandSeats(sd...))
	c, err := rtb.Select(ctx, selector, b...)
	if err != nil {
		xlog.GetWithError(ctx, err).Errorf("invalid request from %s", sup.Name())
		return nil, err
	}
	safe.GoRoutine(ctx, func() {
		for _, s := range sd {
			metrics.Size.With(prometheus.Labels{
				"sup":  sup.Name(),
				"size": s.Size,
				"io":   "in",
			}).Inc()
		}
	})
	return demand.Render(context.Background(), c, req.Id, int(version))
}
