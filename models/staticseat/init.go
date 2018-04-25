package staticseat

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
	staticSeatExp = config.RegisterDuration("crane.models.expire.static.seat", time.Minute, "expire time of static seats")
)

type loader struct {
}

func (loader) Initialize() {
	ctx := context.Background()
	loader := entities.StaticSeatLoader

	staticSeats := pool.NewPool(loader, memorypool.NewMemoryPool(), staticSeatExp.Duration(), 10*time.Second, 3)
	staticSeats.Start(ctx)

	// Wait for the first time load
	<-staticSeats.Notify()

	logrus.Debug("Pool initialized and ready")
}

func init() {
	mysql.Register(&loader{})
}
