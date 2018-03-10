package allads

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/reducer"
	"clickyab.com/crane/demand/rtb"
	"clickyab.com/crane/internal/cyslot"
	"clickyab.com/crane/models/ads"
	"clickyab.com/crane/models/apps"
	"clickyab.com/crane/models/suppliers"
	"clickyab.com/crane/models/website"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/random"
)

var (
	adCTREffect   = config.RegisterInt("crane.rtb.ad_ctr_effect", 50, "ad ctr effect")
	slotCTREffect = config.RegisterInt("crane.rtb.slot_ctr_effect", 50, "slot ctr effect")
)

func AllAdHandler(c context.Context, w http.ResponseWriter, r *http.Request) {

	sup, err := suppliers.GetSupplierByName(r.URL.Query().Get("sup"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	reqType := entity.RequestType(r.URL.Query().Get("type"))
	if !(reqType.IsValid()) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request type"))
		return
	}

	var percentage int64 = 100
	if reqType == entity.RequestTypeNative {
		percentage = 200
	}
	pub := r.URL.Query().Get("pub")
	if !sup.AllowCreate() && pub == "" {
		framework.JSON(w, http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: "publisher required",
		})
		return
	}

	// web or app
	target := r.URL.Query().Get("target")
	brand := r.URL.Query().Get("brand")
	carrier := r.URL.Query().Get("carrier")
	size := r.URL.Query().Get("size")
	latLon := r.URL.Query().Get("latlon")
	os := r.URL.Query().Get("os")
	minBid, _ := strconv.ParseFloat(r.URL.Query().Get("minbid"), 64)

	isp := r.URL.Query().Get("isp") == "t"
	desktop := r.URL.Query().Get("desktop") == "t"
	province := r.URL.Query().Get("province") == "t"
	blacklist := r.URL.Query().Get("blacklist") == "t"
	whitelist := r.URL.Query().Get("whitelist") == "t"

	var cat []string
	if r.URL.Query().Get("cat") != "" {
		cat = strings.Split(r.URL.Query().Get("cat"), ",")
	}

	ua := r.URL.Query().Get("ua")
	if ua == "" {
		ua = r.UserAgent()
	}

	ip := r.URL.Query().Get("ip")
	if ip == "" {
		ip = framework.RealIP(r)
	}

	var (
		publisher entity.Publisher
		selector  []reducer.Filter
	)

	bq := openrtb.BidRequest{}
	if target == "web" {
		publisher, err = website.GetWebSiteOrFake(sup, pub)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("publisher err"))
			return
		}
		bq.Site = &openrtb.Site{
			Inventory: openrtb.Inventory{
				Cat: cat,
			},
		}
		selector = filterWebBuilder(desktop, province, os, isp, whitelist, blacklist, cat)
	} else if target == "app" {
		publisher, err = apps.GetAppOrFake(sup, pub)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("publisher err"))
			return
		}
		bq.App = &openrtb.App{
			Inventory: openrtb.Inventory{
				Cat: cat,
			},
		}
		selector = filterAppBuilder(province, latLon, carrier, brand, isp, whitelist, blacklist, cat)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("target invalid"))
		return
	}

	// filtered ads with errors
	fe := make(map[int64][]string)
	fn := func(id int64, errs []string) {
		fe[id] = errs
	}
	mix := Mix(fn, selector...)

	var ou *openrtb.User
	if latLon != "" {
		var err error
		var lat, lon float64
		if k := strings.Split(latLon, ","); len(k) == 2 {
			lat, err = strconv.ParseFloat(k[0], 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			lon, err = strconv.ParseFloat(k[0], 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			ou = &openrtb.User{
				Geo: &openrtb.Geo{
					Lat: lat,
					Lon: lon,
				},
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("not valid lat/lon"))
			return
		}
	}

	ctr := -1.
	var iSize int
	if size != "" {
		iSize, err = cyslot.GetSize(size)
		if err != nil {
			switch entity.RequestType(target) {
			case entity.RequestTypeNative:
				iSize = 20
			case entity.RequestTypeVast:
				iSize = 9
			default:
				iSize = -1
			}
		}
	}

	bq.User = &openrtb.User{
		ID: <-random.ID,
	}

	ctx, err := builder.NewContext(makeBuilder(carrier, ua, percentage, ip, ou, publisher, sup, bq)...)
	seat := &seat{
		ctr:     ctr,
		size:    iSize,
		rq:      reqType,
		rate:    float64(sup.Rate()),
		minBid:  minBid,
		context: ctx,
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	filteredAds := reducer.Apply(nil, ctx, ads.GetAds(), mix)

	fAds := make([]responseAds, 0)
	for id, ers := range fe {
		a, err := ads.GetAd(id)
		assert.Nil(err)
		fAds = append(fAds, responseAds{
			ID:           a.ID(),
			CampaignID:   a.Campaign().ID(),
			MaxBid:       a.MaxBID(),
			Size:         a.Size(),
			CampaignName: a.Campaign().Name(),
			TargetURL:    a.TargetURL(),
			Assets:       a.Assets(),
			Errors:       ers,
		})
	}

	ex, un := rtb.MinimalSelect(c, ctx, seat, filteredAds)

	framework.JSON(w, http.StatusOK, response{
		FilteredAds:   fAds,
		ExceedAds:     responseBuilder(ex),
		UnderfloorAds: responseBuilder(un),
	})
}

func makeBuilder(carrier, ua string, percentage int64, ip string,
	ou *openrtb.User, publisher entity.Publisher, sup entity.Supplier, bq openrtb.BidRequest) []builder.ShowOptionSetter {
	return []builder.ShowOptionSetter{
		builder.SetTimestamp(),
		builder.SetOSUserAgent(ua),
		builder.SetMinBidPercentage(percentage),
		builder.SetIPLocation(ip, ou, nil),
		builder.SetPublisher(publisher),
		builder.SetStrategy([]string{}, sup),
		builder.SetRate(float64(sup.Rate())),
		builder.SetCappingMode(entity.CappingNone),
		builder.SetUnderfloor(true),
		builder.SetCategory(bq),
		builder.SetTID("", ip, ua),
		builder.SetCarrier(carrier),
		builder.SetCategory(bq),
	}
}

func responseBuilder(ads []entity.SelectedCreative) []responseAds {
	res := make([]responseAds, 0)

	for _, v := range ads {
		x := responseAds{
			ID:            v.ID(),
			CampaignID:    v.Campaign().ID(),
			MaxBid:        v.MaxBID(),
			CampaignName:  v.Campaign().Name(),
			TargetURL:     v.TargetURL(),
			Size:          v.Size(),
			Frequency:     v.Campaign().Frequency(),
			CalculatedCPC: v.CalculatedCPC(),
			CalculatedCPM: v.CalculatedCPM(),
			CalculatedCTR: v.CalculatedCTR(),
			Assets:        v.Assets(),
			Errors:        nil,
		}
		res = append(res, x)
	}

	return res
}
