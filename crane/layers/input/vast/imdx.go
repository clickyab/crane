package vast

import (
	"context"
	"net/http"

	"clickyab.com/crane/crane/builder"
	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/filter"
	"clickyab.com/crane/crane/layers/output/vast"
	"clickyab.com/crane/crane/models"
	"clickyab.com/crane/crane/reducer"
	"clickyab.com/crane/crane/rtb"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/xlog"
)

const vastPath = "/vast/index"

var (
	// TODO : fix this
	vastSelector = reducer.Mix(
		&filter.WebSize{},
		&filter.WebNetwork{},
		&filter.WebMobile{},
		&filter.Desktop{},
		&filter.OS{},
		&filter.WhiteList{},
		&filter.BlackList{},
		&filter.Category{},
		&filter.Province{},
		&filter.ISP{},
	)
)

// vastIndex the vast index route
func vastIndex(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// Hardcoded dear clickyab
	sup, err := models.GetSupplierByName("clickyab")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pubID := r.URL.Query().Get("a")
	pub, err := models.GetWebSiteByPubID(sup, pubID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	us := r.URL.Query().Get("tid")
	parent := r.URL.Query().Get("parent")
	ref := r.URL.Query().Get("ref")
	ip := framework.RealIP(r)
	ua := r.UserAgent()
	b := []builder.ShowOptionSetter{
		builder.SetTimestamp(),
		builder.SetType(entity.RequestTypeVast),
		builder.SetTargetHost(r.URL.Host),
		builder.SetOSUserAgent(ua),
		builder.SetIPLocation(ip),
		builder.SetPublisher(pub),
		builder.SetAlexa(r.UserAgent(), r.Header),
		builder.SetProtocolByRequest(r),
		builder.SetTID(us, ip, ua),

		// Website of demand has no floor cpm and soft floor cpm associated with them
		// TODO : decide about this specific values
		builder.SetFloorCPM(pub.FloorCPM()),
		builder.SetSoftFloorCPM(pub.SoftFloorCPM()),
		builder.SetRate(float64(sup.Rate())),
	}
	// TODO : get this values from config
	b = append(b, builder.SetFloorPercentage(70), builder.SetMinBidPercentage(70))
	b = append(b, builder.SetParent(parent, ref))

	start := r.URL.Query().Get("start") != ""
	mid := r.URL.Query().Get("mid") != ""
	end := r.URL.Query().Get("end") != ""
	l := r.URL.Query().Get("l")
	b = append(b, builder.SetVastSeats(l, pub.PublicID(), start, mid, end))

	c, err := rtb.Select(ctx, vastSelector, b...)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		xlog.GetWithError(ctx, err).Error("invalid request")
		return
	}

	assert.Nil(vast.Render(ctx, w, c))
}
