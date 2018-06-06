package perlocations

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
	//crPerLocationsExp is ads expiration time in redis
	crPerLocationsExp = config.RegisterDuration("crane.models.creatuves.statistics.perlocations", 30*time.Minute, "expire time of crPerLocations. default is 1 hour")
	extraStat         = config.RegisterString("debug.models.creatuves.statistics.perlocations.extra_file", "", "extra file to load for crPerLocations")
)

type pattern struct {
	Data entities.CreativesLocationsReport `json:"data"`
	ID   int64                             `json:"id"`
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
	loader := entities.CRPerLocationsLoader()

	if extraStat.String() != "" {
		loader = pool.DebugLoaderGenerator(loader, extraStat.String(), pattern{})
	}

	crPerLocationsPool = pool.NewPool(loader, memorypool.NewMemoryPool(), crPerLocationsExp.Duration(), 15*time.Second, 3)
	crPerLocationsPool.Start(ctx)

	// Wait for the first time load
	<-crPerLocationsPool.Notify()

	logrus.Debug("crPerLocations pool initialized and ready")
}

func init() {
	mysql.Register(&loader{})
}
