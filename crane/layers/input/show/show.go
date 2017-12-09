package show

import (
	"bytes"
	"context"
	"crypto/md5"
	"net/http"
	"net/url"
	"strconv"

	"clickyab.com/crane/crane/builder"
	"clickyab.com/crane/crane/layers/output/banner"
	"clickyab.com/crane/crane/models"
	"clickyab.com/crane/crane/rtb"
	"clickyab.com/gad/tmp/src/github.com/pkg/errors"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/store/jwt"
)

func extracter(r *url.URL) (map[string]string, error) {
	jt := r.Query().Get("jt")
	if jt == "" {
		return nil, errors.New("jt not found")
	}
	expired, m, err := jwt.NewJWT().Decode([]byte(jt), []string{
		"aid",
		"dom",
		"bid",
		"uaip",
		"susp",
	})
	if err != nil {
		return nil, errors.New("blah bn")
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
		builder.SetOSUserAgent(ua),
		builder.SetIPLocation(ip),
		builder.SetProtocol(r),
		builder.SetQueryParameters(r.URL),
		builder.SetDemandSeats(m["pid"], size),
		builder.SetTID(m["tid"]),
		builder.SetType(m["type"]),
		builder.SetAd(ad),
		builder.SetBid(bid),
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// todo call show worker

	b := &bytes.Buffer{}
	err = banner.Render(ctx, b, c, c.Seats()[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	assert.Nil(w.Write(b.Bytes()))
}
