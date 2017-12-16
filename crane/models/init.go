package models

import (
	"time"

	"context"

	"clickyab.com/crane/crane/models/internal/entities"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/mysql"
	"github.com/clickyab/services/pool"
	"github.com/clickyab/services/pool/drivers/cachepool"
	"github.com/clickyab/services/pool/drivers/memorypool"
)

var (
	websiteExp  = config.RegisterDuration("crane.models.expire.website", time.Hour, "expire time of websites")
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
	websitePubID = pool.NewPool(entities.WebsitePubIDLoader, cachepool.NewCachePool("WP_"), websiteExp.Duration(), 3)
	websitePubID.Start(ctx)

	ads = pool.NewPool(entities.AdLoader, memorypool.NewMemoryPool(), adsExp.Duration(), 3)
	ads.Start(ctx)

	// Wait for the first time load
	<-suppliers.Notify()
	<-websites.Notify()
	<-websitePubID.Notify()
	<-ads.Notify()
}

func init() {
	mysql.Register(&loader{})
}
