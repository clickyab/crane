package distribution

import (
	"fmt"
	"time"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/kv"
)

//Distribution struct type to manage creatives distribution roles
type Distribution struct {
}

const (
	distributeKey = "DS"
)

var (
	distributeExpire = config.RegisterDuration("crane.distribution.expire", 6*time.Hour, "distribution key expiration time. default is 6 hours")
	storeHandler     kv.Kiwi
)

func setDistributeKey(supName, crtype string, publisherID int64, size int, seatID string) {
	storeKey := fmt.Sprintf(
		"%s_%s_%s_%d_%d_%s",
		distributeKey,
		supName,
		crtype,
		publisherID,
		size,
		seatID,
	)

	storeHandler = kv.NewEavStore(storeKey)
}

func getDistributionSelected() map[string]string {
	return storeHandler.AllKeys()
}

func setUpperCreative(ids string) error {
	return storeHandler.SetSubKey("upper", ids).Save(distributeExpire.Duration())
}

func setUnderCreative(ids string) error {
	return storeHandler.SetSubKey("under", ids).Save(distributeExpire.Duration())
}
