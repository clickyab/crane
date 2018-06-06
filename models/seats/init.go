package seats

import (
	"time"

	"context"

	"fmt"

	"clickyab.com/crane/models/internal/entities"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/mysql"
	"github.com/clickyab/services/pool"
	"github.com/clickyab/services/pool/drivers/memorypool"
	"github.com/sirupsen/logrus"
)

var (
	//seatsExp is ads expiration time in redis
	seatsExp  = config.RegisterDuration("crane.models.expire.seats", time.Hour, "expire time of seats. default is 1 hour")
	extraStat = config.RegisterString("debug.models.seats.extra_file", "", "extra file to load for seats")
)

type pattern struct {
	Data entities.Seat `json:"data"`
	ID   int64         `json:"id"`
}

func (p pattern) Value() kv.Serializable {
	return &p.Data
}

func (p pattern) Key() string {
	return fmt.Sprint(p.ID)
}

type loader struct {
}

func (loader) Initialize() {
	ctx := context.Background()
	loader := entities.SeatsLoader()

	if extraStat.String() != "" {
		loader = pool.DebugLoaderGenerator(loader, extraStat.String(), pattern{})
	}

	seatsPool = pool.NewPool(loader, memorypool.NewMemoryPool(), seatsExp.Duration(), 15*time.Second, 3)
	seatsPool.Start(ctx)

	// Wait for the first time load
	<-seatsPool.Notify()

	logrus.Debug("Seats pool initialized and ready")
}

func init() {
	mysql.Register(&loader{})
}
