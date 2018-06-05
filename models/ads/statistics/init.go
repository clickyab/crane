package statistics

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
	//statisticsExp is ads expiration time in redis
	statisticsExp = config.RegisterDuration("crane.models.expire.creatives.statistics", time.Hour, "expire time of network creatives statistics. default is 1 hour")
	extraStat     = config.RegisterString("debug.models.creatives.statistics.extra_file", "", "extra file to load for creatives statistics")
)

type pattern struct {
	Data         entities.CreativeStatistics `json:"data"`
	CreativeType int64                       `json:"ad_type"`
}

func (p pattern) Value() kv.Serializable {
	return &p.Data
}

func (p pattern) Key() string {
	return fmt.Sprint(p.CreativeType)
}

type loader struct {
}

func (loader) Initialize() {
	ctx := context.Background()
	loader := entities.StatisticsLoader

	if extraStat.String() != "" {
		loader = pool.DebugLoaderGenerator(loader, extraStat.String(), pattern{})
	}

	creativesStatisticsPool = pool.NewPool(loader, memorypool.NewMemoryPool(), statisticsExp.Duration(), 15*time.Second, 3)
	creativesStatisticsPool.Start(ctx)

	// Wait for the first time load
	<-creativesStatisticsPool.Notify()

	logrus.Debug("Pool initialized and ready")
}

func init() {
	mysql.Register(&loader{})
}
