package item

import (
	"context"
	"time"

	"github.com/clickyab/services/xlog"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/models/internal/entities"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/mysql"
	"github.com/clickyab/services/pool"
	"github.com/clickyab/services/pool/drivers/memorypool"
	"github.com/sirupsen/logrus"
)

var item pool.Interface

// GetItem return all ads in system
func GetItem(ctx context.Context, s string) entity.Item {
	t, err := item.Get(s, &entities.Asset{})
	if err != nil {
		xlog.GetWithError(ctx, err).Debug()
	}
	return t.(*entities.Asset)
}

// GetItems return all ads in system
func GetItems() map[string]entity.Item {
	data := item.All()
	all := make(map[string]entity.Item)
	for i := range data {
		c := data[i].(entity.Item)
		all[c.Hash()] = c
	}
	return all
}

var (
	adsExp = config.RegisterDuration("crane.models.expire.item", time.Minute*10, "expire time of ads")
)

type loader struct {
}

func (loader) Initialize() {

	item = pool.NewPool(entities.AssetLoader, memorypool.NewMemoryPool(), adsExp.Duration(), 10*time.Second, 3)
	item.Start(context.Background())

	// Wait for the first time load
	<-item.Notify()

	logrus.Debug("Pool initialized and ready")
}

func init() {
	mysql.Register(&loader{})
}
