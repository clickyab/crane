package apps

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
	appExp = config.RegisterDuration("crane.models.expire.app", time.Hour, "expire time of app")
)

type loader struct {
}

func (loader) Initialize() {
	ctx := context.Background()
	app = pool.NewPool(entities.AppLoaderGen(true), cachepool.NewCachePool("APP_"), appExp.Duration(), 30*time.Second, 3)
	app.Start(ctx)

	// Wait for the first time load
	<-app.Notify()

	logrus.Debug("Pool initialized and ready")
}

func init() {
	mysql.Register(&loader{})
}
