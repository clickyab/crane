package pages

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
	//pagesExp is ads expiration time in redis
	pagesExp  = config.RegisterDuration("crane.models.expire.pages", 15*time.Minute, "expire time of pages. default is 1 hour")
	extraStat = config.RegisterString("debug.models.pages.extra_file", "", "extra file to load for pages")
)

type pattern struct {
	Data entities.PublisherPage `json:"data"`
	ID   int64                  `json:"id"`
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
	loader := entities.PagesLoader()

	if extraStat.String() != "" {
		loader = pool.DebugLoaderGenerator(loader, extraStat.String(), pattern{})
	}

	pagesPool = pool.NewPool(loader, memorypool.NewMemoryPool(), pagesExp.Duration(), 15*time.Second, 3)
	pagesPool.Start(ctx)

	// Wait for the first time load
	<-pagesPool.Notify()

	logrus.Debug("Pages pool initialized and ready")
}

func init() {
	mysql.Register(&loader{})
}
