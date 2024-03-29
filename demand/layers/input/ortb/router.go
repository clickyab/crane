package ortb

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/internal/hash"
	"clickyab.com/crane/models/apps"
	"clickyab.com/crane/models/suppliers"
	"clickyab.com/crane/models/website"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
	"github.com/clickyab/services/store/jwt"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
)

type controller struct {
}

type payloadData struct {
	ReserveHash string
	Size        int32
	Type        entity.InputType
	requestType entity.RequestType
	PubType     entity.PublisherType
	TID         string
	Ref         string
	Parent      string
	AdID        int32
	CpID        int32
	CpAdID      int32
	// Ad           entity.Creative
	Supplier     entity.Supplier
	Publisher    entity.Publisher
	Bid          float64
	PublicID     string
	Suspicious   int
	UserAgent    string
	IP           string
	PreviousTime int64
	CPM          float64
	SCPM         float64
	FatFinger    bool
	Tiny         bool
	TV           bool
	CappRegion   string
	CMode        entity.CappingMode // capping mode (none,strict,reset)
	Did          string
	targetURL    string
	cpn          string
}

func extractor(ctx context.Context, r *http.Request) (*payloadData, error) {
	var err error
	jt := xmux.Param(ctx, "jt")
	if jt == "" {
		return nil, errors.New("jt not found")
	}
	pl := payloadData{}
	pl.ReserveHash = xmux.Param(ctx, "rh")
	pl.Type = entity.InputType(xmux.Param(ctx, "type"))
	pl.requestType = entity.RequestType(xmux.Param(ctx, "subtype"))
	pl.TID = r.URL.Query().Get("tid")
	pl.Ref = r.URL.Query().Get("ref")
	pl.Parent = r.URL.Query().Get("parent")

	// check for tv (true view)
	tv := r.URL.Query().Get("tv")
	pl.TV = func(tv string) bool {
		if tv != "" {
			tvInt, err := strconv.ParseInt(tv, 10, 0)
			if err != nil {
				return false
			}
			return tvInt == 1
		}
		return false
	}(tv)

	expired, m, err := jwt.NewJWT().Decode([]byte(jt), "pt", "aid", "sup", "dom", "bid", "uaip",
		"susp", "pid", "now", "cpm", "ff", "t", "cmode", "did", "cpid", "cpadid", "tr", "cpn")
	if err != nil {
		return nil, err
	}

	pl.PreviousTime, err = strconv.ParseInt(m["now"], 10, 0)

	if err != nil {
		return nil, err
	}

	pl.CPM, _ = strconv.ParseFloat(m["cpm"], 64)
	pl.PublicID = m["pid"]
	pl.Did = m["did"]
	// Get the supplier
	pl.Supplier, err = suppliers.GetSupplierByName(m["sup"])
	if err != nil {
		return nil, err
	}

	pl.FatFinger = m["ff"] == "T"
	pl.Tiny = m["t"] == "T"
	pl.SCPM, _ = strconv.ParseFloat(r.URL.Query().Get("scpm"), 64)
	pl.SCPM = pl.SCPM * float64(pl.Supplier.Rate())
	pl.PubType = entity.PublisherType(m["pt"])
	// get the publisher, even its not created then its fine
	if pl.PubType == entity.PublisherTypeWeb {
		pl.Publisher, err = website.GetWebSiteOrFake(pl.Supplier, m["dom"])
	} else if pl.PubType == entity.PublisherTypeApp {
		pl.Publisher, err = apps.GetAppOrFake(pl.Supplier, m["dom"])
	} else {
		err = errors.New("not supported subtype")
	}

	if err != nil {
		return nil, fmt.Errorf("can not find publisher")
	}
	aid, _ := strconv.ParseInt(m["aid"], 10, 64)
	pl.AdID = int32(aid)
	cpid, _ := strconv.ParseInt(m["cpid"], 10, 64)
	pl.CpID = int32(cpid)
	cpadid, _ := strconv.ParseInt(m["cpadid"], 10, 64)
	pl.CpAdID = int32(cpadid)

	pl.targetURL = m["tr"]
	pl.cpn = m["cpn"]
	ts, err := strconv.Atoi(xmux.Param(ctx, "size"))
	if err != nil {
		return nil, fmt.Errorf("invalid size %s", m["size"])
	}
	pl.Size = int32(ts)
	tb, err := strconv.ParseFloat(m["bid"], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid bid %s", m["bid"])
	}
	pl.Bid = tb
	pl.Suspicious, _ = strconv.Atoi(m["susp"])

	pl.UserAgent, pl.IP = r.UserAgent(), framework.RealIP(r)
	mode := 0
	if pl.Publisher.Type() == entity.PublisherTypeApp {
		mode = 1
	}
	tmp := hash.Sign(mode, pl.ReserveHash, fmt.Sprint(pl.Size), string(pl.Type), pl.UserAgent, pl.IP)
	if tmp != m["uaip"] {
		pl.Suspicious = 98
	}

	if expired {
		pl.Suspicious = 99
	}
	if pl.Supplier.Name() == "chavoosh" || pl.Supplier.Name() == "chavoosh-ne2" {
		pl.Suspicious = 0
	}
	pl.CappRegion = r.URL.Query().Get("reg")
	cMode, err := strconv.Atoi(m["cmode"])
	assert.Nil(err)
	pl.CMode = entity.CappingMode(cMode)

	return &pl, nil
}

func (controller) Routes(m framework.Mux) {
	m.GET("banner", showPath, showBanner)
	m.GET("click", clickPath, clickBanner)
	m.GET("balance", balancePath, balancer)
	m.GET("pixel", pixelPath, showPixel)
	m.POST("ortb", demandPath, openRTBInput)
	m.GET("notice", noticePath, noticeHandler)
	// Backward compatible route
	m.RootMux().GET(conversionPath, xhandler.HandlerFuncC(conversionHandler))
	m.RootMux().GET(conversionTrackerJS, xhandler.HandlerFuncC(trackerJS))
	// New routes
	m.GET("conversion", conversionPath, conversionHandler)
	m.GET("conversion-js", conversionTrackerJS, trackerJS)
	m.GET("capping", cappingURL, applyCapp)
}

func init() {
	router.Register(controller{})
}
