package entities

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"net/url"
	"unicode/utf8"

	"github.com/clickyab/services/assert"

	"github.com/clickyab/services/config"

	"github.com/clickyab/services/xlog"

	"github.com/clickyab/services/simplehash"

	"clickyab.com/crane/demand/entity"
	openrtb "clickyab.com/crane/openrtb/v2.5"

	"github.com/sirupsen/logrus"

	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/mysql"
)

// Asset for targeting
type Asset struct {
	FID        int32            `db:"id"`
	FHash      string           `db:"hash"`
	FURL       string           `db:"url"`
	FSKU       string           `db:"sku"`
	FBrand     string           `db:"brand"`
	FImg       string           `db:"img"`
	FTitle     string           `db:"title"`
	FPrice     mysql.NullInt64  `db:"price"`
	FDiscount  mysql.NullInt64  `db:"discount"`
	FMime      string           `db:"mime"`
	FCat       mysql.NullString `db:"cat"`
	FAvailable int              `db:"is_available"`
	campaign   entity.Campaign  `db:"-"`
	capping    entity.Capping   `db:"-"`
	assets     []entity.Asset   `db:"-"`
}

// Available return true if product is available in advertiser inventory
func (a *Asset) Available() bool {
	if a.FAvailable == 1 {
		return true
	}
	return false
}

// SetCampaign just set campaign ;)
func (a *Asset) SetCampaign(c entity.Campaign) {
	a.campaign = c
}

// Type return the type of ad
func (a *Asset) Type() entity.AdType {
	return entity.AdTypeNative
}

// Campaign return the type of ad
func (a *Asset) Campaign() entity.Campaign {
	return a.campaign
}

// AdCTR the ad ctr from database (its not calculated from )
func (a *Asset) AdCTR() float32 {
	return float32(a.campaign.CTR())
}

// Target return the target of this campaign
func (a *Asset) Target() entity.Target {
	return entity.TargetNative
}

// Size returns ads size
func (a *Asset) Size() int32 {
	return 20
}

// Width return the width
func (a *Asset) Width() int32 {
	return 0
}

// Height return the height of banner
func (a *Asset) Height() int32 {
	return 0
}

// Duration of the ad if it have meaning. normally usable for vast, in second
func (a *Asset) Duration() int32 {
	return 0
}

// Capping return the current capping object
func (a *Asset) Capping() entity.Capping {
	return a.capping
}

// SetCapping set the current capping
func (a *Asset) SetCapping(c entity.Capping) {
	a.capping = c
}

// Attributes return the ad specific attributes
func (a *Asset) Attributes() map[string]interface{} {
	return map[string]interface{}{}
}

// Media return image of ad
func (a *Asset) Media() string {
	return a.FImg
}

// TargetURL return ad target of the target
func (a *Asset) TargetURL() string {
	return a.FURL
}

// CampaignAdID return campaign_ad primary
func (a *Asset) CampaignAdID() int32 {
	return a.campaign.ID() * -1
}

// MimeType of media
func (a *Asset) MimeType() string {
	return a.FMime
}

// {
// "banner_title_text_type":"فانوس وب | طراحی سایت، سئو، طراحی فروشگاه اینترنتی",
// "product":"http:\/\/static.clickyab.com\/ad\/product_53426_27136_1549981654.jpeg",
// "banner_description_text_type":"",
// "link_text_type":"https:\/\/fanoosweb.ir",
// "w":495,
// "h":400
// }

// Asset return the asset that pass all filters and type is exactly matched the value
func (a *Asset) Asset(assetType entity.AssetType, sub int, filter ...entity.AssetFilter) []entity.Asset {

	var res []entity.Asset
	// Ignore if the assets is empty
	if len(a.assets) == 0 {
		return res
	}
bigLoop:
	for i := range a.assets {
		if a.assets[i].Type != assetType || a.assets[i].SubType != sub {
			continue
		}
		for j := range filter {
			if !filter[j](&a.assets[i]) {
				continue bigLoop
			}
		}
		res = append(res, a.assets[i])
	}
	return res
}

// Assets return the assets
func (a *Asset) Assets() []entity.Asset {
	return a.assets
}

// Cat of product
func (a *Asset) Cat() []openrtb.ContentCategory {
	return a.campaign.Category()
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
	if a.FPrice.Valid {
		return a.FPrice.Int64
	}
	return -1
}

// Discount of product if any
func (a *Asset) Discount() int64 {
	if a.FDiscount.Valid {
		return a.FDiscount.Int64
	}
	return -1
}

// ID of item
func (a *Asset) ID() int32 {
	return a.FID + 1e6*-1
}

// Hash of url
func (a *Asset) Hash() string {
	return a.FHash
}

// Encode is the encode of this function
func (a *Asset) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(a)
}

// Decode is the decode function
func (a *Asset) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(a)
}

var imgconv = config.RegisterString("crane.demand.cropper", "http://t.clickyab.com/cropper/", "")

// AssetLoader for caching advertiser product
func AssetLoader(_ context.Context) (map[string]kv.Serializable, error) {

	var res []Asset
	q := fmt.Sprintf(`select id,hash,url,sku,brand,img,title,price,discount,is_available,cat from list_asset`)
	_, err := NewManager().GetRDbMap().Select(&res, q)
	if err != nil {
		return nil, err
	}

	for i := range res {
		img, err := url.Parse(imgconv.String())
		assert.Nil(err)
		qs := img.Query()

		qs.Add("w", "500")
		qs.Add("x", "500")
		qs.Add("h", "315")
		qs.Add("y", "315")
		qs.Add("url", res[i].FImg)
		img.RawQuery = qs.Encode()
		res[i].assets = []entity.Asset{{
			MimeType: "image/jpeg",
			Type:     entity.AssetTypeImage,
			SubType:  entity.AssetTypeImageSubTypeMain,
			Width:    int32(500),
			Height:   int32(315),
			Data:     img.String(),
		}, {
			MimeType: "text/html",
			Type:     entity.AssetTypeText,
			SubType:  entity.AssetTypeTextSubTypeTitle,
			Len:      int32(utf8.RuneCountInString(res[i].FTitle)),
			Data:     res[i].FTitle,
		}}
	}

	k := make(map[string]kv.Serializable)
	for i := range res {
		k[res[i].FHash] = &res[i]
	}
	logrus.Debugf("Load %d items", len(k))
	return k, nil
}

// AddAssets add multiple product asset
func AddAssets(ctx context.Context, a []entity.Item) error {
	q := "INSERT INTO list_asset (hash,url,sku,brand,img,title,discount,is_available) VALUES "
	p := "(?,?,?,?,?,?,?,?),"
	x := " ON DUPLICATE KEY UPDATE hash=VALUES(hash)"
	params := make([]interface{}, 0)

	for _, e := range a {
		a := func() int32 {
			if e.Available() {
				return 1
			}
			return 0
		}()
		q += p
		params = append(params, simplehash.SHA1(e.URL()), e.URL(), e.SKU(), e.Brand(), e.Image(), e.Title(), e.Discount(), a)
	}
	q = q[:len(q)-1] + x

	_, err := NewManager().GetWDbMap().Exec(q, params...)
	if err != nil {
		xlog.GetWithError(ctx, err).Debug(q, params)
	}
	return err
}
