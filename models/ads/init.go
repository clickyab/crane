package ads

import (
	"sync/atomic"
	"time"

	"clickyab.com/crane/models/ads/statistics/locationctr"

	"github.com/clickyab/services/safe"

	"context"

	"fmt"

	"clickyab.com/crane/models/internal/entities"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/mysql"
	"github.com/clickyab/services/pool"
	"github.com/clickyab/services/pool/drivers/memorypool"
	"github.com/sirupsen/logrus"
)

var (
	//AdsExp is ads expiration time in redis
	AdsExp   = config.RegisterDuration("crane.models.expire.ads", time.Minute, "expire time of ads")
	extraAds = config.RegisterString("debug.models.ads.extra_file", "", "extra file to load for ads")
	started  int64
)

type pattern struct {
	Data entities.Advertise `json:"data"`
	ID   int64              `json:"id"`
}

func (p pattern) Value() kv.Serializable {
	return &p.Data
}

func (p pattern) Key() string {
	return fmt.Sprint(p.ID)
}

type loader struct {
}

func (loader) Initialize() {
	ctx := context.Background()
	loader := entities.AdLoader

	if extraAds.String() != "" {
		loader = pool.DebugLoaderGenerator(loader, extraAds.String(), pattern{})
	}

	ads = pool.NewPool(loader, memorypool.NewMemoryPool(), AdsExp.Duration(), 10*time.Second, 3)
	ads.Start(ctx)

	// Wait for the first time load
	<-ads.Notify()
	logrus.Debug("Pool of creatives initialized and ready")
	listenToUpdates()
}

func listenToUpdates() {
	if !atomic.CompareAndSwapInt64(&started, 0, 1) {
		return
	}

	ctx := context.Background()
	safe.ContinuesGoRoutine(ctx, func(x context.CancelFunc) time.Duration {
		for i := 0; ; i++ {
			select {
			case <-ads.Notify():
				err := locationctr.Load(getIds())
				if err != nil {
					logrus.Warn(err)
				} else {
					logrus.Debug("Pool of creatives ctr per page and location initialized and ready")
				}
			default:
			}
		}

		return AdsExp.Duration()
	})
}

func init() {
	mysql.Register(&loader{})
}
