package ortb

import (
	"errors"

	"clickyab.com/crane/openrtb"
)

//
//import (
//	"encoding/json"
//	"errors"
//	"fmt"
//	"strings"
//
//	"clickyab.com/crane/demand/builder"
//	"clickyab.com/crane/demand/entity"
//	"clickyab.com/crane/models/suppliers"
//	"clickyab.com/crane/openrtb"
//	"github.com/clickyab/services/assert"
//	"github.com/clickyab/services/xlog"
//	"golang.org/x/net/context"
//)
//
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

//type server struct {
//}

//
//func (*server) Ortb(ctx context.Context, req *openrtb.BidRequest) (*openrtb.BidResponse, error) {
//	res := &openrtb.BidResponse{}
//	res.Id = req.GetId()
//	token := req.GetToken()
//	sup, err := suppliers.GetSupplierByToken(token)
//	if err != nil {
//		e := fmt.Errorf("supplier with token %s not found", token)
//		xlog.GetWithError(ctx, err).Debug(e)
//		return nil, e
//	}
//
//	var (
//		tiny, fatFinger, prevent, underfloor bool
//		strategy                             []string
//		capping                              = entity.CappingStrict
//	)
//	if req.Ext != nil {
//		req.GetExt()
//		fatFinger = req.Ext.GetFatFinger()
//		prevent = req.Ext.GetPrevent()
//		underfloor = req.Ext.GetUnderfloor()
//		tiny = req.Ext.GetTiny()
//		strategy = req.Ext.GetStrategy()
//		if req.Ext.GetCapping() == openrtb.Capping_Reset {
//			capping = entity.CappingReset
//		}
//	}
//
//	if err := validate(req); err != nil {
//		return res, err
//	}
//
//	var domain, bundle string
//	if req.GetSite() != nil {
//		domain = req.GetSite().GetDomain()
//	}
//	if req.GetApp() != nil {
//		bundle = req.GetApp().GetBundle()
//	}
//
//	pub, selector, ps, prevent, err := handlePublisherSelector(domain, bundle, sup, prevent)
//	if err != nil {
//		return res, err
//	}
//	proto := entity.HTTP
//	for _, m := range req.GetImp() {
//		if m.GetSecure() {
//			proto = entity.HTTPS
//			break
//		}
//	}
//
//	ua := ""
//	ip := ""
//	if req.GetDevice() != nil {
//		ua = strings.Trim(req.GetDevice().GetUa(), "\n\t ")
//		ip = strings.Trim(req.GetDevice().GetIp(), "\n\t ")
//	}
//	us := ""
//	if req.GetUser() != nil {
//		us = req.GetUser().GetId()
//	}
//
//	if ua == "" || ip == "" {
//		err := fmt.Errorf("no ip/no ua")
//		xlog.GetWithError(ctx, err).Debug("invalid request")
//		return res, err
//	}
//
//	b := []builder.ShowOptionSetter{
//		builder.SetTimestamp(),
//		builder.SetTargetHost(sup.ShowDomain()),
//		builder.SetOSUserAgent(ua),
//		builder.SetIPLocation(ip, req.GetUser(), req.GetDevice()),
//		builder.SetPublisher(pub),
//		builder.SetProtocol(proto),
//		builder.SetTID(us, ip, ua),
//		builder.SetNoTiny(!tiny),
//		builder.SetBannerMarkup(sup),
//		builder.SetFatFinger(fatFinger),
//		builder.SetStrategy(strategy, sup),
//		builder.SetRate(float64(sup.Rate())),
//		builder.SetPreventDefault(prevent),
//		builder.SetCappingMode(capping),
//		builder.SetUnderfloor(underfloor),
//		builder.SetCategory(req),
//	}
//	// TODO : if we need to implement native/app/vast then the next line must be activated and customized
//	//b = append(b, builder.SetFloorPercentage(100), builder.SetMinBidPercentage(100))
//
//	b = setGRPCPublisherCustomContext(*req, b)
//
//	return res, nil
//}
//
//func (*server) OrtbStream(req openrtb.OrtbService_OrtbStreamServer) error {
//	panic("implement me")
//}
//
//func grpcSeatDetail(req openrtb.BidRequest) ([]builder.DemandSeatData, bool) {
//	var (
//		imp    = req.Imp
//		seats  = make([]builder.DemandSeatData, 0)
//		w, h   int32
//		vast   bool
//		assets []*openrtb.NativeRequest_Asset
//	)
//	for _, m := range imp {
//		var t entity.RequestType
//		if m.GetVideo() != nil {
//			w, h = m.GetVideo().GetW(), m.GetVideo().GetH()
//			t = entity.RequestTypeVast
//			// We just support version 3
//			var v3 bool
//			for _, v := range m.GetVideo().GetProtocols() {
//				if v == openrtb.Protocol_VAST_3_0 {
//					v3 = true
//					break
//				}
//			}
//			if !v3 {
//				continue
//			}
//			vast = true
//		} else if m.GetBanner() != nil {
//			w, h = m.GetBanner().W, m.GetBanner().H
//			t = entity.RequestTypeBanner
//		} else if m.GetNative() != nil {
//			t = entity.RequestTypeNative
//			req := &openrtb.NativeRequest{}
//			err := json.Unmarshal([]byte(m.GetNative().GetRequest()), &req)
//			assert.Nil(err)
//			assets = req.Assets
//		}
//		var (
//			ext = make(simpleMap)
//		)
//
//		seats = append(seats, builder.DemandSeatData{
//			MinBid: m.GetBidfloor(),
//			PubID:  m.GetId(),
//			Size:   fmt.Sprintf("%dx%d", w, h),
//			Type:   t,
//			Video:  m.GetVideo(),
//			Banner: m.GetBanner(),
//			Assets: assets,
//			MinCPC: ext.Float64("min_cpc"),
//		})
//	}
//	return seats, vast
//}
//
//func setGRPCPublisherCustomContext(payload openrtb.BidRequest, b []builder.ShowOptionSetter) []builder.ShowOptionSetter {
//	if payload.GetSite() != nil {
//		b = append(b, builder.SetParent(payload.GetSite().GetPage(), payload.GetSite().GetRef()))
//	}
//	if payload.GetApp() != nil {
//		b = append(b, builder.SetConnType(payload.GetDevice().GetConnectiontype()))
//		b = append(b, builder.SetCarrier(payload.Device.Carrier))
//		b = append(b, builder.SetBrand(payload.Device.Model))
//	}
//	return b
//}
