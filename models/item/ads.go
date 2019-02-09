package item

import (
	"context"
	"time"

	openrtb "clickyab.com/crane/openrtb/v2.5"
	"github.com/clickyab/services/simplehash"

	"github.com/clickyab/services/xlog"

	"github.com/sirupsen/logrus"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/models/internal/entities"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/mysql"
	"github.com/clickyab/services/pool"
	"github.com/clickyab/services/pool/drivers/memorypool"
)

var item pool.Interface

// GetItem return all ads in system
func GetItem(ctx context.Context, s string) entity.Item {
	t, err := item.Get(s, &entities.Asset{})
	if err != nil {
		xlog.GetWithError(ctx, err).Debug("GET ITEM")
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
	adsExp = config.RegisterDuration("crane.models.expire.item", time.Minute*1, "expire time of ads")
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

// Asset for product
type Asset struct {
	FList       string        `json:"list"`
	FURL        string        `json:"url"`
	FImg        string        `json:"img"`
	FTitle      string        `json:"title"`
	FPrice      int64         `json:"price"`
	FDiscount   int64         `json:"discount"`
	FSKU        string        `json:"sku"`
	IsAvailable bool          `json:"is_available"`
	FCategory   []string      `json:"category"`
	FBrand      string        `json:"brand"`
	FAvailable  bool          `json:"is_available"`
	User        *openrtb.User `json:"-"`
}

// Available for availability
func (a *Asset) Available() bool {
	return a.FAvailable
}

// ID of item
func (a *Asset) ID() int32 {
	return 0
}

// Hash of url
func (a *Asset) Hash() string {
	return simplehash.SHA1(a.FURL)
}

// URL of product
func (a *Asset) URL() string {
	return a.FURL
}

// SKU is advertiser product id
func (a *Asset) SKU() string {
	return a.FSKU
}

// Brand of product
func (a *Asset) Brand() string {
	return a.FBrand
}

// Image url of product
func (a *Asset) Image() string {
	return a.FImg
}

// Title of product
func (a *Asset) Title() string {
	return a.FTitle
}

// Price of product
func (a *Asset) Price() int64 {
	return a.FPrice
}

// Discount of product if any
func (a *Asset) Discount() int64 {
	return a.FDiscount
}

// Cat of product
func (a *Asset) Cat() []openrtb.ContentCategory {
	return nil
}

// AddAssets for adding for products
func AddAssets(ctx context.Context, as []entity.Item) error {
	return entities.AddAssets(ctx, as)
}

// CheckList to be sure target list exists
func CheckList(s string) (*entities.List, error) {
	return entities.CheckList(s)
}
