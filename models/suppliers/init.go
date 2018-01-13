package suppliers

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
	supplierExp = config.RegisterDuration("crane.models.expire.supplier", time.Hour, "expire time of supplier")
)

type loader struct {
}

func (loader) Initialize() {
	ctx := context.Background()
	suppliers = pool.NewPool(entities.SupplierLoader, memorypool.NewMemoryPool(), supplierExp.Duration(), 10*time.Second, 3)
	suppliers.Start(ctx)
	suppliersByName = pool.NewPool(entities.SupplierLoaderByName, memorypool.NewMemoryPool(), supplierExp.Duration(), 10*time.Second, 3)
	suppliersByName.Start(ctx)

	// Wait for the first time load
	<-suppliers.Notify()
	<-suppliersByName.Notify()

	logrus.Debug("Pool initialized and ready")
}

func init() {
	mysql.Register(&loader{})
}
