package ortb

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
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
	lock          = sync.RWMutex{}
	regionPath    = map[string]string{}
)

type configInit struct{}

func (configInit) Initialize() config.DescriptiveLayer {
	return nil
}

func (configInit) Loaded() {
	lock.Lock()
	defer lock.Unlock()

	for _, v := range strings.Split(regions.String(), ",") {
		x := strings.Split(v, ":")
		if len(x) != 2 {
			continue
		}
		regionPath[x[0]] = x[1]
	}
}

func getRegion(key string) string {
	lock.RLock()
	defer lock.RUnlock()
	return regionPath[key]
}

func setCapping(ctx context.Context, pl *payloadData, proto string) {
	if (pl.CappRegion == currentRegion.String() || pl.CappRegion == "") && pl.CMode != entity.CappingNone {
		capping.StoreCapping(pl.CMode, pl.TID, pl.Ad.ID())
		return
	}

	rURL := getRegion(pl.CappRegion)
	if rURL == "" {
		xlog.GetWithError(ctx, fmt.Errorf("invalid region : %s", pl.CappRegion)).Error("invalid region")
		return
	}

	var httpClient = &http.Client{}

	urlPath := router.MustPath("capping", map[string]string{
		"ad_id":     fmt.Sprint(pl.Ad.ID()),
		"user_id":   pl.TID,
		"capp_mode": fmt.Sprint(pl.CMode),
	})

	cappUpdateURL := &url.URL{
		Host:   rURL,
		Scheme: proto,
		Path:   urlPath,
	}

	cappRequest, err := http.NewRequest("GET", cappUpdateURL.String(), nil)
	if err != nil {
		xlog.GetWithError(ctx, err).Debug("request create failed")
		return
	}

	err = callCappingUpdate(ctx, httpClient, cappRequest)
	if err != nil {
		xlog.GetWithError(ctx, err).Error("call capping on the other region failed")
	}
}

func callCappingUpdate(ctx context.Context, httpClient *http.Client, r *http.Request) error {
	resp, err := httpClient.Do(r.WithContext(ctx))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status %d", resp.StatusCode)
	}
	return nil
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

	// Build context
	c, err := builder.NewContext(
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
		builder.SetFullSeats(pl.PublicID, pl.Size, pl.ReserveHash, pl.Ad, pl.Bid, time.Now().Unix(), pl.CPM, pl.SCPM, pl.requestType),
	)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exp, cnl := context.WithTimeout(ctx, 10*time.Second)
	safe.GoRoutine(exp, func() {
		job := show.NewImpressionJob(c, c.Seats()...)
		broker.Publish(job)
		cnl()
	})

	safe.GoRoutine(ctx, func() {
		setCapping(ctx, pl, c.Protocol().String())
	})

	assert.Nil(banner.Render(ctx, w, c))
}

func init() {
	config.Register(&configInit{})
}
