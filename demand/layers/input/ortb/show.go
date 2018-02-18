package ortb

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/capping"
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/layers/output/banner"
	"clickyab.com/crane/workers/show"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework/router"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/safe"
	"github.com/clickyab/services/xlog"
	"github.com/sirupsen/logrus"
)

const showPath = "/banner/:rh/:size/:type/:subtype/:jt"

var (
	regions       = config.RegisterString("crane.regions.urls", "us:us-demand.com,fr:fr-demand.com", "determine valid regions and related domain for update capp url")
	currentRegion = config.RegisterString("crane.regions.current", "fr", "determine current region")
)

func setCapping(ctx context.Context, pl *payloadData, u string, proto string) {

	if (pl.CappRegion == currentRegion.String() || pl.CappRegion == "") && pl.CMode != entity.CappingNone {
		capping.StoreCapping(pl.CMode, u, pl.Ad.ID())
		return
	}

	for _, v := range strings.Split(regions.String(), ",") {
		if x := strings.Split(v, ":"); x[0] == pl.CappRegion {
			var httpClient = &http.Client{}

			urlPath := router.MustPath("capping", map[string]string{
				"adID":     fmt.Sprint(pl.Ad.ID()),
				"userID":   u,
				"cappMode": fmt.Sprint(pl.CMode),
			})

			cappUpdateURL := &url.URL{
				Host:   x[1],
				Scheme: proto,
				Path:   urlPath,
			}

			cappRequest, err := http.NewRequest("GET", cappUpdateURL.String(), nil)
			if err != nil {
				xlog.GetWithError(ctx, err).Debug("request create failed")
				return
			}

			callCappingUpdate(ctx, httpClient, cappRequest)
			return
		}
	}

	xlog.GetWithFields(ctx, logrus.Fields{"capping_mode": pl.CMode,
		"payload":        pl,
		"region":         regions.String(),
		"current_region": currentRegion.String()}).
		Error("invalid region for capping")

}

// show is handler for show ad requestType
func showBanner(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pl, err := extractor(ctx, r)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	counter := kv.NewAEAVStore(pl.ReserveHash, clickExpire.Duration()+time.Hour).IncSubKey("I", 1)
	if counter > 1 {
		// Duplicate impression!
		pl.Suspicious = 3
	}
	b := []builder.ShowOptionSetter{
		builder.SetTimestamp(),
		builder.SetOSUserAgent(pl.UserAgent),
		builder.SetNoTiny(!pl.Tiny),
		builder.SetCappingMode(entity.CappingMode(pl.CMode)),
		builder.SetTargetHost(r.Host),
		builder.SetIPLocation(pl.IP, nil, nil),
		builder.SetProtocolByRequest(r),
		builder.SetParent(pl.Parent, pl.Ref),
		builder.SetTID(pl.TID, pl.IP, pl.UserAgent),
		builder.SetPublisher(pl.Publisher),
		builder.SetSuspicious(pl.Suspicious),
		builder.SetFatFinger(pl.FatFinger),
	}

	b = append(b, builder.SetFullSeats(pl.PublicID, pl.Size, pl.ReserveHash, pl.Ad, pl.Bid, time.Now().Unix(), pl.CPM, pl.SCPM, pl.requestType))
	// Build context
	c, err := builder.NewContext(b...)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exp, _ := context.WithTimeout(ctx, 10*time.Second)
	safe.GoRoutine(exp, func() {
		job := show.NewImpressionJob(c, c.Seats()...)
		broker.Publish(job)
	})

	safe.GoRoutine(ctx, func() {
		setCapping(ctx, pl, c.User().ID(), c.Protocol().String())
	})

	assert.Nil(banner.Render(ctx, w, c))
}

func callCappingUpdate(ctx context.Context, httpClient *http.Client, r *http.Request) {
	resp, err := httpClient.Do(r.WithContext(ctx))
	if err != nil {
		xlog.GetWithError(ctx, err).Debug("request do failed")
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("invalid status %d", resp.StatusCode)
		xlog.GetWithError(ctx, err).Debug("request do status failed")
		return
	}
}
