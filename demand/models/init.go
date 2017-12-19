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
	mobileExp   = config.RegisterDuration("crane.models.expire.mobile", time.Hour, "expire time of networks,brand and carrier")
)

type loader struct {
}

func (loader) Initialize() {
	ctx := context.Background()
	networks = pool.NewPool(entities.NetworkLoader, memorypool.NewMemoryPool(), mobileExp.Duration(), 3)
	networks.Start(ctx)
	carriers = pool.NewPool(entities.CarrierLoader, memorypool.NewMemoryPool(), mobileExp.Duration(), 3)
	carriers.Start(ctx)
	brands = pool.NewPool(entities.BrandLoader, memorypool.NewMemoryPool(), mobileExp.Duration(), 3)
	brands.Start(ctx)
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
	<-networks.Notify()
	<-carriers.Notify()
	<-brands.Notify()
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
