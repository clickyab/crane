package ads

import (
	"time"

	"context"

	"clickyab.com/crane/models/internal/entities"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/mysql"
	"github.com/clickyab/services/pool"
	"github.com/clickyab/services/pool/drivers/memorypool"
	"github.com/sirupsen/logrus"
)

var (
	adsExp = config.RegisterDuration("crane.models.expire.ads", time.Minute, "expire time of ads")
)

type loader struct {
}

func (loader) Initialize() {
	ctx := context.Background()
	ads = pool.NewPool(entities.AdLoader, memorypool.NewMemoryPool(), adsExp.Duration(), 3)
	ads.Start(ctx)

	// Wait for the first time load
	<-ads.Notify()

	logrus.Debug("Pool initialized and ready")
}

func init() {
	mysql.Register(&loader{})
}
