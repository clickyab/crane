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
	"clickyab.com/crane/models/apps"
	"clickyab.com/crane/supplier/client"
	"clickyab.com/crane/supplier/layer/output"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/random"
	"github.com/mssola/user_agent"
)

var (
	sup    entity.Supplier = &supplier{}
	server                 = config.RegisterString("crane.supplier.banner.url", "", "route for app")
	method                 = config.RegisterString("crane.supplier.app.method", "POST", "method for app request")
)

func getApp(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pub, err := apps.GetApp(sup, r.URL.Query().Get("package"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	bs, width, height, full, err := size(r.URL.Query().Get("adsMedia"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	q := &openrtb.BidRequest{
		ID:      <-random.ID,
		AllImps: 1,
		Imp: []openrtb.Impression{
			{
				ID:     <-random.ID,
				Secure: secure(r),
				Banner: &openrtb.Banner{
					W:  width,
					H:  height,
					ID: <-random.ID,
				},
			},
		},
	}
	ext := map[string]interface{}{
		"cid":          r.URL.Query().Get("cid"),
		"lac":          r.URL.Query().Get("lac"),
		"capping_mode": "reset",
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if output.RenderApp(ctx, w, res, full, sdkVers, bs) != nil {
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
	data, _ := strconv.Atoi(r.URL.Query().Get("data"))
	h, _ := strconv.Atoi(r.URL.Query().Get("screenHeight"))
	w, _ := strconv.Atoi(r.URL.Query().Get("screenWidth"))
	ppi, _ := strconv.Atoi(r.URL.Query().Get("screenDensity"))
	lat, lon := 0.0, 0.0
	gps := strings.Split(r.URL.Query().Get("gps"), ",")

	if len(gps) == 2 {
		lat, _ = strconv.ParseFloat(gps[0], 64)
		lon, _ = strconv.ParseFloat(gps[1], 64)
	}

	o.App = &openrtb.App{
		Bundle: r.URL.Query().Get("package"),
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
		ConnType: data,
		Geo: &openrtb.Geo{
			Lat: lat,
			Lon: lon,
		},
	}

	o.User = &openrtb.User{
		ID: r.URL.Query().Get("androidid"),
		Geo: &openrtb.Geo{
			Lat: lat,
			Lon: lon,
		},
	}

}

func secure(r *http.Request) openrtb.NumberOrString {
	if framework.Scheme(r) == "https" {
		return openrtb.NumberOrString(1)
	}
	return openrtb.NumberOrString(0)
}
