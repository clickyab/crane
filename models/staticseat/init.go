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
	staticSeatExp = config.RegisterDuration("crane.models.expire.static.seat", 2*time.Minute, "expire time of static seats")
)

type loader struct {
}

func (loader) Initialize() {
	ctx := context.Background()
	loader := entities.StaticSeatLoader

	seats = pool.NewPool(loader, memorypool.NewMemoryPool(), staticSeatExp.Duration(), 10*time.Second, 3)
	seats.Start(ctx)

	// Wait for the first time load
	<-seats.Notify()

	logrus.Debug("Pool initialized and ready")
}

func init() {
	mysql.Register(&loader{})
}
