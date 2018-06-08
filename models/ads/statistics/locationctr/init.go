package locationctr

import (
	"sync/atomic"
	"time"

	"clickyab.com/crane/models/internal/entities"

	"fmt"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/mysql"
	"github.com/clickyab/services/pool/drivers/cachepool"
)

var (
	//crlocationctrExp is ads expiration time in redis
	crlocationctrExp = config.RegisterDuration("crane.models.creatuves.statistics.locationctr", 3*time.Hour, "expire time of crlocationctr. default is 1 hour")
	started          int64
)

type loader struct {
}

func (loader) Initialize() {
}

func makeConnections(ids []int64) {
	for i := range ids {
		if cn := locationCTRPools[ids[i]]; cn.key == "" {
			k := fmt.Sprintf("CTR_PER_PAGE_%d", ids[i])

			locationCTRPools[ids[i]] = pools{
				key:     k,
				driver:  cachepool.NewCachePool(k),
				scanner: kv.NewScanner(k),
			}
		}
	}
}

//Load pool data
func Load(ids []int64) error {
	if !atomic.CompareAndSwapInt64(&started, 0, 1) {
		return nil
	}

	makeConnections(ids)

	var finalIds []int64

	for i := range ids {
		_, check := locationCTRPools[ids[i]].scanner.Next(1)
		if !check {
			finalIds = append(finalIds, ids[i])
		}
	}

	if len(finalIds) == 0 {
		atomic.SwapInt64(&started, 0)
		return nil
	}

	newData, err := entities.CRlocationctrLoader(finalIds)
	if err != nil {
		atomic.SwapInt64(&started, 0)
		return err
	}

	for i := range newData {
		err := locationCTRPools[ids[i]].driver.Store(newData[i], crlocationctrExp.Duration())
		if err != nil {
			atomic.SwapInt64(&started, 0)
			return err
		}
	}

	atomic.SwapInt64(&started, 0)
	return nil
}

func init() {
	mysql.Register(&loader{})
}
