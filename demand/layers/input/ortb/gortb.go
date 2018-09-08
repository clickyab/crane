package ortb

import (
	"errors"
	"fmt"
	"strings"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/models/suppliers"
	"clickyab.com/crane/openrtb"
	"github.com/clickyab/services/xlog"
	"golang.org/x/net/context"
)

// Validation errors
var (
	ErrInvalidReqNoID        = errors.New("openrtb: request ID missing")
	ErrInvalidReqNoImps      = errors.New("openrtb: request has no impressions")
	ErrInvalidReqMultiInv    = errors.New("openrtb: request has multiple inventory sources") // has site and app
	ErrInvalidImpNoID        = errors.New("openrtb: impression ID missing")
	ErrInvalidImpMultiAssets = errors.New("openrtb: impression has multiple assets") // at least two out of Banner, Video, Native

	ErrInvalidVideoNoMimes       = errors.New("openrtb: video has no mimes")
	ErrInvalidVideoNoLinearity   = errors.New("openrtb: video linearity missing")
	ErrInvalidVideoNoMinDuration = errors.New("openrtb: video min-duration missing")
	ErrInvalidVideoNoMaxDuration = errors.New("openrtb: video max-duration missing")
	ErrInvalidVideoNoProtocols   = errors.New("openrtb: video protocols missing")
)

type server struct {
}

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

	} else if req.GetApp() == nil && req.GetSite() == nil {
		return ErrInvalidReqMultiInv
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

func (*server) Ortb(ctx context.Context, req *openrtb.BidRequest) (*openrtb.BidResponse, error) {
	res := &openrtb.BidResponse{}
	res.Id = req.GetId()
	token := req.GetToken()
	sup, err := suppliers.GetSupplierByToken(token)
	if err != nil {
		e := fmt.Errorf("supplier with token %s not found", token)
		xlog.GetWithError(ctx, err).Debug(e)
		return nil, e
	}

	var (
		tiny, fatFinger, prevent, underfloor bool
		strategy                             []string
		capping                              = entity.CappingStrict
	)
	if req.Ext != nil {
		req.GetExt()
		fatFinger = req.Ext.GetFatFinger()
		prevent = req.Ext.GetPrevent()
		underfloor = req.Ext.GetUnderfloor()
		tiny = req.Ext.GetTiny()
		strategy = req.Ext.GetStrategy()
		if req.Ext.GetCapping() == openrtb.Capping_Reset {
			capping = entity.CappingReset
		}
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
		return res, err
	}
	proto := entity.HTTP
	for _, m := range req.GetImp() {
		if m.GetSecure() {
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
	us := ""
	if req.GetUser() != nil {
		us = req.GetUser().GetId()
	}

	if ua == "" || ip == "" {
		err := fmt.Errorf("no ip/no ua")
		xlog.GetWithError(ctx, err).Debug("invalid request")
		return res, err
	}

	b := []builder.ShowOptionSetter{
		builder.SetTimestamp(),
		builder.SetTargetHost(sup.ShowDomain()),
		builder.SetOSUserAgent(ua),
		builder.SetGRPCIPLocation(ip, req.GetUser(), req.GetDevice()),
		builder.SetPublisher(pub),
		builder.SetProtocol(proto),
		builder.SetTID(us, ip, ua),
		builder.SetNoTiny(!tiny),
		builder.SetBannerMarkup(sup),
		builder.SetFatFinger(fatFinger),
		builder.SetStrategy(strategy, sup),
		builder.SetRate(float64(sup.Rate())),
		builder.SetPreventDefault(prevent),
		builder.SetCappingMode(capping),
		builder.SetUnderfloor(underfloor),
		builder.SetGRPCCategory(req),
	}
	// TODO : if we need to implement native/app/vast then the next line must be activated and customized
	//b = append(b, builder.SetFloorPercentage(100), builder.SetMinBidPercentage(100))

	return res, nil
}

func (*server) OrtbStream(req openrtb.OrtbService_OrtbStreamServer) error {
	panic("implement me")
}

func setGRPCPublisherCustomContext(payload openrtb.BidRequest, b []builder.ShowOptionSetter) []builder.ShowOptionSetter {
	if payload.GetSite() != nil {
		b = append(b, builder.SetParent(payload.GetSite().GetPage(), payload.GetSite().GetRef()))
	}
	if payload.GetApp() != nil {
		b = append(b, builder.SetConnType(int(payload.GetDevice().GetConnectiontype())))
		b = append(b, builder.SetCarrier(payload.Device.Carrier))
		b = append(b, builder.SetBrand(payload.Device.Model))
	}
	return b
}
