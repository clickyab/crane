package show

import (
	"context"
	"crypto/md5"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"clickyab.com/crane/crane/builder"
	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/layers/output/banner"
	"clickyab.com/crane/crane/models"
	"clickyab.com/crane/crane/rtb"
	"clickyab.com/crane/crane/workers/show"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/safe"
	"github.com/clickyab/services/store/jwt"
	"github.com/rs/xmux"
)

const showPath = "/banner/:rh/:size/:type/:jt"

func extractor(ctx context.Context, r *url.URL) (map[string]string, error) {
	jt := xmux.Param(ctx, "jt")
	if jt == "" {
		return nil, errors.New("jt not found")
	}
	expired, m, err := jwt.NewJWT().Decode([]byte(jt), "aid", "sup", "dom", "bid", "uaip", "susp", "pid")
	if err != nil {
		return nil, err
	}
	if expired {
		m["susp"] = "99"
	}
	m["rh"] = xmux.Param(ctx, "rh")
	m["size"] = xmux.Param(ctx, "size")
	m["type"] = xmux.Param(ctx, "type")
	m["tid"] = r.Query().Get("tid")
	m["ref"] = r.Query().Get("ref")
	m["parent"] = r.Query().Get("parent")
	return m, nil
}

// show is handler for show ad request
func showBanner(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	m, err := extractor(ctx, r.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ua, ip := r.UserAgent(), framework.RealIP(r)
	if string(md5.New().Sum([]byte(r.UserAgent()+framework.RealIP(r)))) != m["uaip"] {
		m["susp"] = "98"
	}
	// Get the supplier
	sup, err := models.GetSupplierByName(m["sup"])
	assert.Nil(err)
	// get the publisher, even its not created then its fine
	pub, err := models.GetWebSite(sup, m["dom"])
	assert.Nil(err)
	aid, err := strconv.ParseInt(m["aid"], 10, 64)
	assert.Nil(err)
	ad, err := models.GetAd(aid)
	assert.Nil(err)
	size, err := strconv.Atoi(m["size"])
	assert.Nil(err)
	bid, err := strconv.ParseFloat(m["bid"], 64)
	assert.Nil(err)
	susp, err := strconv.Atoi(m["susp"])
	assert.Nil(err)

	// Build context
	c, err := rtb.Select(ctx, nil,
		builder.SetTimestamp(),
		builder.SetOSUserAgent(ua),
		builder.SetRequest(r.Host, r.Method),
		builder.SetIPLocation(ip),
		builder.SetAlexa(ua, http.Header{}),
		builder.SetProtocolByRequest(r),
		builder.SetParent(m["parent"], m["ref"]),
		builder.SetTID(m["tid"]),
		builder.SetType(entity.RequestType(m["type"])),
		builder.SetPublisher(pub),
		builder.SetSuspicious(susp),
		builder.SetFullSeats(m["pid"], size, m["rh"], ad, bid),
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exp, _ := context.WithTimeout(ctx, 10*time.Second)
	safe.GoRoutine(exp, func() {
		job := show.NewImpressionJob(c, c.Seats()...)
		broker.Publish(job)
	})

	assert.Nil(banner.Render(ctx, w, c))
}
