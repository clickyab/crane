package website

import (
	"time"

	"context"

	"clickyab.com/crane/models/internal/entities"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/mysql"
	"github.com/clickyab/services/pool"
	"github.com/clickyab/services/pool/drivers/cachepool"
	"github.com/sirupsen/logrus"
)

var (
	websiteExp = config.RegisterDuration("crane.models.expire.website", time.Hour, "expire time of websites")
)

type loader struct {
}

func (loader) Initialize() {
	ctx := context.Background()

	websites = pool.NewPool(entities.WebsiteLoader, cachepool.NewCachePool("WS_"), websiteExp.Duration(), 10*time.Second, 3)
	websites.Start(ctx)

	// Wait for the first time load
	<-websites.Notify()

	logrus.Debug("Pool initialized and ready")
}

func init() {
	mysql.Register(&loader{})
}
