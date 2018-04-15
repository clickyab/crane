package app

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"encoding/json"

	"errors"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/models/clickyabapps"
	"clickyab.com/crane/supplier/client"
	"clickyab.com/crane/supplier/layers/internal/supplier"
	"clickyab.com/crane/supplier/layers/output"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/simplehash"
	"github.com/clickyab/services/xlog"
	"github.com/mssola/user_agent"
)

var (
	sup             = supplier.NewClickyab()
	server          = config.RegisterString("crane.supplier.banner.url", "", "route for app")
	method          = config.RegisterString("crane.supplier.app.method", "POST", "method for app request")
	clickyabNetwork = map[string]networkConn{
		"2G":   cellular2G,
		"EDGE": cellular2G,
		"GPRS": cellular2G,
		"3G":   cellular3G,
		"4G":   cellular4G,
	}
)

type networkConn int

const (
	unknownNetwork = 0

	cellular2G = 4
	cellular3G = 5
	cellular4G = 6
)

func getApp(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pub, err := clickyabapps.GetApp(sup, r.URL.Query().Get("token"))
	if err != nil {
		xlog.GetWithError(ctx, err).Debugf("token invalid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bs, width, height, full, err := size(r.URL.Query().Get("adsMedia"))
	if err != nil {
		xlog.GetWithError(ctx, err).Debugf("adsMedia invalid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// calculate min cpc and insert in impression ext
	impExt := map[string]interface{}{
		"min_cpc": pub.MinCPC(string(entity.RequestTypeBanner)),
	}
	iExt, err := json.Marshal(impExt)
	assert.Nil(err)

	q := &openrtb.BidRequest{
		ID: <-random.ID,
		App: &openrtb.App{
			Bundle: pub.Name(),
		},
		Imp: []openrtb.Impression{
			{
				ID:     <-random.ID,
				Secure: secure(r),
				Banner: &openrtb.Banner{
					W:  width,
					H:  height,
					ID: <-random.ID,
				},
				BidFloor: float64(pub.FloorCPM()),
				Ext:      iExt,
			},
		},
	}
	ext := map[string]interface{}{
		"cid":          r.URL.Query().Get("cid"),
		"lac":          r.URL.Query().Get("lac"),
		"capping_mode": "reset",
		"underfloor":   true,
		"tiny_mark":    false,
	}
	if _, ok := pub.Attributes()[entity.PAFatFinger]; ok {
		ext["fat_finger"] = true
	}
	sdkVers, _ := strconv.ParseInt(r.URL.Query().Get("clickyabVersion"), 10, 0)
	if sdkVers <= 4 {
		// older version of sdk (pre 5) use a method to handle click which is not correct.
		// this is a workaround for that
		ext["prevent_default"] = true
	}

	j, err := json.Marshal(ext)
	assert.Nil(err)
	q.Ext = j

	allData(r, q)

	res, err := client.Call(ctx, method.String(), server.String(), q)
	if err != nil {
		xlog.GetWithError(ctx, err).Debugf("call failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if output.RenderApp(ctx, w, res, full, sdkVers, bs) != nil {
		xlog.GetWithError(ctx, err).Debugf("render failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func size(s string) (int, int, int, string, error) {
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
	lat, lon := 0.0, 0.0
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
		IP:       framework.RealIP(r),
		OS:       ua.OS(),
		DNT:      dnt,
		Carrier:  r.URL.Query().Get("carrier"),
		Language: r.URL.Query().Get("lang"),
		H:        h,
		W:        w,
		MCCMNC:   fmt.Sprintf("%s-%s", r.URL.Query().Get("mcc"), r.URL.Query().Get("mnc")),
		PPI:      ppi,
		Model:    r.URL.Query().Get("brand"),
		HwVer:    r.URL.Query().Get("model"),
		OSVer:    r.URL.Query().Get("androidVersion"),
		UA:       r.UserAgent(),
		ConnType: int(getConnType(network)),
		Geo: &openrtb.Geo{
			Lat: lat,
			Lon: lon,
		},
	}

	o.User = &openrtb.User{
		ID: appUserIDGenerator(androidID, deviceID, operator, model),
		Geo: &openrtb.Geo{
			Lat: lat,
			Lon: lon,
		},
	}

}

// getConnType convert clickyab network to openrtb
func getConnType(network string) networkConn {
	val, ok := clickyabNetwork[network]
	if ok {
		return val
	}
	return unknownNetwork
}

// appUserIDGenerator create cop id for app
func appUserIDGenerator(androidID, deviceID, operator, model string) string {
	return simplehash.MD5(fmt.Sprintf("%s%s%s%s", androidID, deviceID, operator, model))
}

// secure check openrtb protocol (http/https)
func secure(r *http.Request) openrtb.NumberOrString {
	if framework.Scheme(r) == "https" {
		return openrtb.NumberOrString(1)
	}
	return openrtb.NumberOrString(0)
}
