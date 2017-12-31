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
	suppliers = pool.NewPool(entities.SupplierLoader, memorypool.NewMemoryPool(), supplierExp.Duration(), 3)
	suppliers.Start(ctx)
	suppliersByName = pool.NewPool(entities.SupplierLoaderByName, memorypool.NewMemoryPool(), supplierExp.Duration(), 3)
	suppliersByName.Start(ctx)

	// Wait for the first time load
	res := fanIn(suppliers.Notify(), suppliersByName.Notify())
	for i := 0; i < 2; i++ {
		<-res
	}

	logrus.Debug("Pool initialized and ready")
}

func fanIn(input1, input2 <-chan time.Time) <-chan time.Time {
	c := make(chan time.Time)
	go func() {
		for {
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			}
		}
	}()

	return c
}

func init() {
	mysql.Register(&loader{})
}
