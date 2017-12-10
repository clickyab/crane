package show

import (
	"bytes"
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
)

func extracter(r *url.URL) (map[string]string, error) {
	jt := r.Query().Get("jt")
	if jt == "" {
		return nil, errors.New("jt not found")
	}
	expired, m, err := jwt.NewJWT().Decode([]byte(jt), "aid", "dom", "bid", "uaip", "susp")
	if err != nil {
		return nil, err
	}
	if expired {
		if _, ok := m["susp"]; ok {
			m["clickexpired"] = "true"
		} else {
			m["showexpired"] = "true"
		}
	}
	m["rh"] = r.Query().Get("rh")
	m["tid"] = r.Query().Get("tid")
	m["ref"] = r.Query().Get("resf")
	m["parent"] = r.Query().Get("parent")
	m["size"] = r.Query().Get("size")
	return m, nil
}

// Show is handler for show ad request
func Show(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	m, err := extracter(r.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ua, ip := r.UserAgent(), framework.RealIP(r)
	if string(md5.New().Sum([]byte(r.UserAgent()+framework.RealIP(r)))) != m["uaip"] {
		m["susp"] = "mismatch user agent or ip"
	}
	aid, err := strconv.ParseInt(m["aid"], 10, 64)
	assert.Nil(err)
	ad, err := models.GetAd(aid)
	assert.Nil(err)
	size, err := strconv.Atoi(m["size"])
	assert.Nil(err)
	bid, err := strconv.ParseFloat(m["bid"], 64)
	assert.Nil(err)

	c, err := rtb.Select(ctx, nil,
		builder.SetTimestamp(),
		builder.SetOSUserAgent(ua),
		builder.SetIPLocation(ip),
		builder.SetProtocol(r),
		builder.SetQueryParameters(r.URL),
		builder.SetDemandSeats([]builder.SeatDetail{{PubID: m["pid"], Size: size}}),
		builder.SetTID(m["tid"]),
		builder.SetType(entity.RequestType(m["type"])),
		builder.SetAd(ad),
		builder.SetBid(bid),
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pub, _ := context.WithTimeout(ctx, 10*time.Second)
	safe.GoRoutine(pub, func() {
		job := show.NewImpressionJob(c, c.Seats()[0])
		broker.Publish(job)
	})

	b := &bytes.Buffer{}
	err = banner.Render(ctx, b, c, c.Seats()[0])
	assert.Nil(err)

	w.WriteHeader(http.StatusOK)
	assert.Nil(w.Write(b.Bytes()))
}
