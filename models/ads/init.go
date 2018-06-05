package ads

import (
	"time"

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

	logrus.Debug("Pool of creatives statistics initialized and ready")
}

func init() {
	mysql.Register(&loader{})
}
