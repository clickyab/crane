package suppliers

import (
	"time"

	"context"

	"clickyab.com/crane/models/internal/entities"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/mysql"
	"github.com/clickyab/services/pool"
	"github.com/clickyab/services/pool/drivers/memorypool"
	"github.com/sirupsen/logrus"
)

var (
	supplierExp = config.RegisterDuration("crane.models.expire.supplier", 1*time.Hour, "expire time of supplier")
	extraSup    = config.RegisterString("debug.models.supplier.extra_file", "", "extra file to load for suppliers")
)

type tokenPattern struct {
	Data  entities.Supplier `json:"data"`
	Token string            `json:"token"`
}

func (p tokenPattern) Value() kv.Serializable {
	return &p.Data
}

func (p tokenPattern) Key() string {
	return p.Token
}

type namePattern struct {
	Data  entities.Supplier `json:"data"`
	Token string            `json:"token"`
}

func (p namePattern) Value() kv.Serializable {
	return &p.Data
}

func (p namePattern) Key() string {
	return p.Data.FName
}

type loader struct {
}

func (loader) Initialize() {
	ctx := context.Background()
	tokenLoader := entities.SupplierLoader
	if extraSup.String() != "" {
		tokenLoader = pool.DebugLoaderGenerator(tokenLoader, extraSup.String(), tokenPattern{})
	}
	suppliers = pool.NewPool(tokenLoader, memorypool.NewMemoryPool(), supplierExp.Duration(), 10*time.Second, 3)
	suppliers.Start(ctx)

	nameLoader := entities.SupplierLoaderByName
	if extraSup.String() != "" {
		nameLoader = pool.DebugLoaderGenerator(nameLoader, extraSup.String(), namePattern{})
	}
	suppliersByName = pool.NewPool(nameLoader, memorypool.NewMemoryPool(), supplierExp.Duration(), 10*time.Second, 3)
	suppliersByName.Start(ctx)

	// Wait for the first time load
	<-suppliers.Notify()
	<-suppliersByName.Notify()

	logrus.Debug("Pool initialized and ready")
}

func init() {
	mysql.Register(&loader{})
}
