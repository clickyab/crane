package ortb

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"strings"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/workers/click"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/safe"
)

const (
	clickPath = "/click/:rh/:size/:type/:subtype/:jt"
	prefix    = "CLK_LIMIT"
	format    = "20060102"
)

var (
	clickExpire     = config.RegisterDuration("crane.context.seat.click_exp", 72*time.Hour, "determine how long click url is valid")
	dailyClickLimit = config.RegisterInt64("crane.context.seat.click_limit", 3, "determine limit click for ip per day")
)

// clickBanner is handler for click ad requestType
func clickBanner(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pl, err := extractor(ctx, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// This is a very important thing in click expiry detection. so be aware of it
	//TODO : if we lose redis somehow, it can lead to a problematic duplicate click,
	//TODO : create an offline job to check duplicate click hash in the past 72 hours
	counter := kv.NewAEAVStore(pl.ReserveHash, clickExpire.Duration()+time.Hour).IncSubKey("C", 1)
	if counter > 1 {
		// Duplicate click!
		pl.Suspicious = 1
	}
	perDay := kv.NewAEAVStore(fmt.Sprintf("%s_%s_%s", prefix, time.Now().Format(format), pl.IP), 24*time.Hour).IncSubKey("C", 1)
	if perDay > dailyClickLimit.Int64() {
		pl.Suspicious = 96
	}

	// Build context
	c, err := builder.NewContext(
		builder.SetTimestamp(),
		builder.SetOSUserAgent(pl.UserAgent),
		builder.SetTargetHost(r.Host),
		builder.SetIPLocation(pl.IP),
		builder.SetProtocolByRequest(r),
		builder.SetTID(pl.TID, pl.IP, pl.UserAgent),
		builder.SetPublisher(pl.Publisher),
		builder.SetSuspicious(pl.Suspicious),
		builder.SetFatFinger(pl.FatFinger),
		builder.SetFullSeats(pl.PublicID, pl.Size, pl.ReserveHash, pl.Ad, pl.Bid, pl.PreviousTime, pl.CPM, pl.SCPM, pl.requestType),
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exp, _ := context.WithTimeout(ctx, 10*time.Second)
	safe.GoRoutine(exp, func() {
		job := click.NewClickJob(c)
		broker.Publish(job)
	})
	body := replaceParameters(pl.Ad.TargetURL(), pl.Publisher.Name(), pl.Ad.Campaign().Name(), pl.ReserveHash, pl.IP)

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
