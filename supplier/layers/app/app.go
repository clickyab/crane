package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/models/clickyabapps"
	"clickyab.com/crane/openrtb/v2.5"
	"clickyab.com/crane/supplier/client"
	"clickyab.com/crane/supplier/layers/internal/supplier"
	"clickyab.com/crane/supplier/layers/output"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/simplehash"
	"github.com/clickyab/services/xlog"
	"github.com/mssola/user_agent"
)

var (
	//server          = config.RegisterString("crane.supplier.banner.url", "", "route for app")
	sup             = supplier.NewClickyab()
	clickyabNetwork = map[string]openrtb.ConnectionType{
		"2G":   openrtb.ConnectionType_CELLULAT_NETWORK_2G,
		"EDGE": openrtb.ConnectionType_CELLULAT_NETWORK_2G,
		"GPRS": openrtb.ConnectionType_CELLULAT_NETWORK_2G,
		"3G":   openrtb.ConnectionType_CELLULAT_NETWORK_3G,
		"4G":   openrtb.ConnectionType_CELLULAT_NETWORK_4G,
	}
)

func getApp(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pub, err := clickyabapps.GetApp(sup, r.URL.Query().Get("token"))
	if err != nil {
		xlog.GetWithError(ctx, err).Debugf("token invalid")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("cly-error", "invalid token")
		return
	}
	bs, width, height, full, err := size(r.URL.Query().Get("adsMedia"))
	if err != nil {
		xlog.GetWithError(ctx, err).Debugf("adsMedia invalid")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("cly-error", "adsMedia invalid")

		return
	}

	q := &openrtb.BidRequest{
		Id: <-random.ID,
		DistributionchannelOneof: &openrtb.BidRequest_App{
			App: &openrtb.App{
				Bundle: pub.Name(),
			},
		},
		Imp: []*openrtb.Imp{
			{
				Id: <-random.ID,
				Secure: func() int32 {
					if framework.Scheme(r) == "https" {
						return 1
					}
					return 0
				}(),
				Banner: &openrtb.Banner{
					W:  width,
					H:  height,
					Id: <-random.ID,
				},
				Bidfloor: float64(pub.FloorCPM()),
				Ext: &openrtb.Imp_Ext{
					Mincpc: pub.MinCPC(string(entity.RequestTypeBanner)),
				},
			},
		},
	}
	sdkVers, _ := strconv.ParseInt(r.URL.Query().Get("clickyabVersion"), 10, 0)
	ext := &openrtb.BidRequest_Ext{
		Capping:    openrtb.Capping_Reset,
		Underfloor: true,
		Tiny:       false,
		Cid:        r.URL.Query().Get("cid"),
		Lac:        r.URL.Query().Get("lac"),
		Prevent:    sdkVers <= 4,
	}

	if _, ok := pub.Attributes()[entity.PAFatFinger]; ok {
		ext.FatFinger = true
	}
	q.Ext = ext

	allData(r, q)

	//res, err := client.Call(ctx, server.String(), q)
	res, err := client.StreamCall(ctx, q)
	if err != nil {
		xlog.GetWithError(ctx, err).Debugf("call failed")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("cly-error", "call failed")

		return
	}

	if output.RenderApp(ctx, w, res, full, sdkVers, bs) != nil {
		xlog.GetWithError(ctx, err).Debugf("render failed")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("cly-error", "render failed")

		return
	}

}

func size(s string) (int32, int32, int32, string, error) {
	switch strings.ToLower(s) {
	case "banner":
		return 8, 320, 50, "", nil
	case "largebanner":
		return 3, 300, 250, "", nil
	case "xlargebannerportrait":
		return 16, 320, 480, "", nil
	case "fullbannerportrait":
		return 16, 320, 480, "portrait", nil
	case "xlargebannerlandscap":
		return 17, 48, 320, "", nil
	case "fullbannerlandscape":
		return 17, 48, 320, "landscape", nil
	default:
		return 0, 0, 0, "", errors.New("not valid size")
	}
}

func allData(r *http.Request, o *openrtb.BidRequest) {

	ua := user_agent.New(r.UserAgent())
	dnt, _ := strconv.Atoi(r.Header.Get("DNT"))
	network := r.URL.Query().Get("network") //2g,3g,4g,...
	h, _ := strconv.Atoi(r.URL.Query().Get("screenHeight"))
	w, _ := strconv.Atoi(r.URL.Query().Get("screenWidth"))
	ppi, _ := strconv.Atoi(r.URL.Query().Get("screenDensity"))
	var lat, lon = 0.0, 0.0
	gps := strings.Split(r.URL.Query().Get("gps"), ",")

	// 4 data for creating user id
	androidID := r.URL.Query().Get("androidid")
	deviceID := r.URL.Query().Get("deviceid")
	operator := r.URL.Query().Get("operator")
	model := r.URL.Query().Get("model") // samsung,huawei,...

	if len(gps) == 2 {
		lat, _ = strconv.ParseFloat(gps[0], 64)
		lon, _ = strconv.ParseFloat(gps[1], 64)
	}

	o.Device = &openrtb.Device{
		Ip:             framework.RealIP(r),
		Os:             ua.OS(),
		Dnt:            int32(dnt),
		Carrier:        r.URL.Query().Get("carrier"),
		Language:       r.URL.Query().Get("lang"),
		H:              int32(h),
		W:              int32(w),
		Mccmnc:         fmt.Sprintf("%s-%s", r.URL.Query().Get("mcc"), r.URL.Query().Get("mnc")),
		Ppi:            int32(ppi),
		Model:          r.URL.Query().Get("brand"),
		Hwv:            r.URL.Query().Get("model"),
		Osv:            r.URL.Query().Get("androidVersion"),
		Ua:             r.UserAgent(),
		Connectiontype: getConnType(network),
		Geo: &openrtb.Geo{
			Lat: float32(lat),
			Lon: float32(lon),
		},
	}

	o.User = &openrtb.User{
		Id: appUserIDGenerator(androidID, deviceID, operator, model),
		Geo: &openrtb.Geo{
			Lat: float32(lat),
			Lon: float32(lon),
		},
	}

}

// getConnType convert clickyab network to openrtb
func getConnType(network string) openrtb.ConnectionType {
	val, ok := clickyabNetwork[network]
	if ok {
		return val
	}
	return openrtb.ConnectionType_CELLULAT_NETWORK_UNKNOWN
}

// appUserIDGenerator create cop id for app
func appUserIDGenerator(androidID, deviceID, operator, model string) string {
	return simplehash.MD5(fmt.Sprintf("%s%s%s%s", androidID, deviceID, operator, model))
}
