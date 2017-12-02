package models

import (
	"time"

	"clickyab.com/crane/crane/models/internal/entities"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/mysql"
	"github.com/clickyab/services/pool"
	"github.com/clickyab/services/pool/drivers/cachepool"
	"github.com/clickyab/services/pool/drivers/memorypool"
)

var (
	websiteExp = config.RegisterDuration("crane.models.expire.website", time.Hour, "expire time of websites")
	brandExp   = config.RegisterDuration("crane.models.expire.brand", time.Hour*24, "expire time of brands")
	networkExp = config.RegisterDuration("crane.models.expire.network", time.Hour*24, "expire time of network")
	carrierExp = config.RegisterDuration("crane.models.expire.carrier", time.Hour*24, "expire time of carrier")
	adsExp     = config.RegisterDuration("crane.models.expire.ads", time.Minute, "expire time of ads")
)

type loader struct {
}

func (loader) Initialize() {
	websites = pool.NewPool(entities.WebsiteLoader, cachepool.NewCachePool("WS_"), websiteExp.Duration(), 3)
	brands = pool.NewPool(entities.BrandLoader, memorypool.NewMemoryPool(), brandExp.Duration(), 3)
	networks = pool.NewPool(entities.NetworkLoader, memorypool.NewMemoryPool(), networkExp.Duration(), 3)
	carriers = pool.NewPool(entities.CarrierLoader, memorypool.NewMemoryPool(), carrierExp.Duration(), 3)
	ads = pool.NewPool(entities.AdLoader, memorypool.NewMemoryPool(), adsExp.Duration(), 3)

	// Wait for the first time load
	<-websites.Notify()
	<-brands.Notify()
	<-networks.Notify()
	<-carriers.Notify()
	<-ads.Notify()
}

func init() {
	mysql.Register(&loader{})
}
