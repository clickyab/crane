package entities

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"

	openrtb "clickyab.com/crane/openrtb/v2.5"
	"github.com/sirupsen/logrus"

	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/mysql"
)

type asset struct {
	FID       int32            `db:"id"`
	FHash     string           `db:"hash"`
	FURL      string           `db:"url"`
	FSKU      string           `db:"sku"`
	FBrand    string           `db:"brand"`
	FImg      string           `db:"img"`
	FTitle    string           `db:"title"`
	FPrice    int32            `db:"price"`
	FDiscount int32            `db:"discount"`
	FCat      mysql.NullString `db:"cat"`
}

func (a *asset) Cat() []openrtb.ContentCategory {
	panic("implement me")
}

func (a *asset) URL() string {
	return a.FURL
}

func (a *asset) SKU() string {
	return a.FSKU
}

func (a *asset) Brand() string {
	return a.FBrand
}

func (a *asset) Image() string {
	return a.FImg
}

func (a *asset) Title() string {
	return a.FTitle
}

func (a *asset) Price() int32 {
	return a.FPrice
}

func (a *asset) Discount() int32 {
	return a.FDiscount
}

func (a *asset) ID() int32 {
	return a.FID
}

func (a *asset) Hash() string {
	return a.FHash
}

func (a *asset) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(a)
}

func (a *asset) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(a)
}

// AssetLoader for caching advertiser product
func AssetLoader(_ context.Context) (map[string]kv.Serializable, error) {

	var res []asset
	q := fmt.Sprintf(`select id,hash,url,sku,brand,img,title,price,discount,is_available,cat from list_asset where is_available = 1`)
	_, err := NewManager().GetRDbMap().Select(&res, q)
	if err != nil {
		return nil, err
	}

	k := make(map[string]kv.Serializable)
	for i := range res {
		k[res[i].FHash] = &res[i]
	}
	logrus.Debugf("Load %d items", len(k))
	return k, nil
}

func init() {
	_, _ = AssetLoader(context.Background())
}
