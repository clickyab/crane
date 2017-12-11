package common

import (
	"context"
	"net/http"
	"time"

	"strings"

	"clickyab.com/crane/crane/builder"
	"clickyab.com/crane/crane/workers/click"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/safe"
)

const clickPath = "/click/:rh/:size/:type/:jt"

// clickBanner is handler for click ad request
func clickBanner(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pl, err := extractor(ctx, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Build context
	c, err := builder.NewContext(
		builder.SetTimestamp(),
		builder.SetOSUserAgent(pl.UserAgent),
		builder.SetRequest(r.Host, r.Method),
		builder.SetIPLocation(pl.IP),
		builder.SetAlexa(pl.UserAgent, http.Header{}),
		builder.SetProtocolByRequest(r),
		builder.SetTID(pl.TID),
		builder.SetType(pl.Type),
		builder.SetPublisher(pl.Publisher),
		builder.SetSuspicious(pl.Suspicious),
		builder.SetFullSeats(pl.PublicID, pl.Size, pl.ReserveHash, pl.Ad, pl.Bid),
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exp, _ := context.WithTimeout(ctx, 10*time.Second)
	safe.GoRoutine(exp, func() {
		job := show.NewClickJob(c)
		broker.Publish(job)
	})
	body := replaceParameters(pl.Ad.AdTarget(), pl.Publisher.Name(), pl.Ad.Campaign().Name(), pl.ReserveHash, pl.IP)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(body))
	assert.Nil(err)
}

func replaceParameters(url, domain, campaign, impID, ip string) string {
	r := strings.NewReplacer(
		"[app]",
		domain,
		"[domain]",
		domain,
		"[campaign]",
		campaign,
		"{app}",
		domain,
		"{domain}",
		domain,
		"{campaign}",
		campaign,
		"{imp_id}",
		impID,
		"{ip}",
		ip,
		"[ip]",
		ip,
	)

	url = r.Replace(url)
	return `<html><head><title>` + url + `</title><meta name="robots" content="nofollow"/></head>
			<body><script>window.setTimeout( function() { window.location.href = '` + url + `' }, 500 );</script></body>
			</html>`
}
