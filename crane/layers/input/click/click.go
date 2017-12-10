package click

import (
	"context"
	"crypto/md5"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"fmt"

	"clickyab.com/crane/crane/builder"
	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/models"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/safe"
	"github.com/clickyab/services/store/jwt"
	"github.com/rs/xmux"
	"github.com/sirupsen/logrus"
)

const clickPath = "/click/:rh/:size/:type/:jt"

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
		m["susp"] = "98"
	}
	m["rh"] = xmux.Param(ctx, "rh")
	m["size"] = xmux.Param(ctx, "size")
	m["type"] = xmux.Param(ctx, "type")
	m["tid"] = r.Query().Get("tid")
	m["ref"] = r.Query().Get("ref")
	m["parent"] = r.Query().Get("parent")
	return m, nil
}

// clickBanner is handler for click ad request
func clickBanner(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	m, err := extractor(ctx, r.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//get ad
	aid, err := strconv.ParseInt(m["aid"], 10, 64)
	assert.Nil(err)
	ad, err := models.GetAd(aid)
	assert.Nil(err)
	size, err := strconv.Atoi(m["size"])
	assert.Nil(err)
	bid, err := strconv.ParseFloat(m["bid"], 64)
	assert.Nil(err)

	ua, ip := r.UserAgent(), framework.RealIP(r)

	//check for fraud
	if string(md5.New().Sum([]byte(m["rh"]+fmt.Sprintf("%d", size)+m["type"]+ua+ip))) != m["uaip"] {
		m["susp"] = "98"
	}
	susp, err := strconv.Atoi(m["susp"])
	assert.Nil(err)
	sup, err := models.GetSupplierByName(m["sup"])
	assert.Nil(err)
	pub, err := models.GetWebSite(sup, m["dom"])
	assert.Nil(err)
	// Build context
	c, err := builder.NewContext(
		builder.SetTimestamp(),
		builder.SetOSUserAgent(ua),
		builder.SetRequest(r.Host, r.Method),
		builder.SetIPLocation(ip),
		builder.SetAlexa(ua, http.Header{}),
		builder.SetProtocolByRequest(r),
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
		// call web worker
		logrus.Debug(c)
	})
}
