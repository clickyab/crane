package models

import (
	"time"

	"context"

	"clickyab.com/crane/demand/models/internal/entities"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/mysql"
	"github.com/clickyab/services/pool"
	"github.com/clickyab/services/pool/drivers/cachepool"
	"github.com/clickyab/services/pool/drivers/memorypool"
	"github.com/sirupsen/logrus"
)

var (
	websiteExp  = config.RegisterDuration("crane.models.expire.website", time.Hour, "expire time of websites")
	appExp      = config.RegisterDuration("crane.models.expire.app", time.Hour, "expire time of apps")
	supplierExp = config.RegisterDuration("crane.models.expire.supplier", time.Hour, "expire time of supplier")
	adsExp      = config.RegisterDuration("crane.models.expire.ads", time.Minute, "expire time of ads")
)

type loader struct {
}

func (loader) Initialize() {
	ctx := context.Background()

	suppliers = pool.NewPool(entities.SupplierLoader, memorypool.NewMemoryPool(), supplierExp.Duration(), 3)
	suppliers.Start(ctx)
	suppliersByName = pool.NewPool(entities.SupplierLoaderByName, memorypool.NewMemoryPool(), supplierExp.Duration(), 3)
	suppliersByName.Start(ctx)

	websites = pool.NewPool(entities.WebsiteLoader, cachepool.NewCachePool("WS_"), websiteExp.Duration(), 3)
	websites.Start(ctx)
	apps = pool.NewPool(entities.AppLoader, cachepool.NewCachePool("APP_"), appExp.Duration(), 3)
	apps.Start(ctx)
	ads = pool.NewPool(entities.AdLoader, memorypool.NewMemoryPool(), adsExp.Duration(), 3)
	ads.Start(ctx)

	// Wait for the first time load
	<-suppliers.Notify()
	<-apps.Notify()
	<-suppliersByName.Notify()
	<-websites.Notify()
	<-ads.Notify()

	logrus.Debug("Pool initialized and ready")
}

func init() {
	mysql.Register(&loader{})
}
