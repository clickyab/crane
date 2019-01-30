package entities

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"unicode/utf8"

	"clickyab.com/crane/demand/entity"
	openrtb "clickyab.com/crane/openrtb/v2.5"

	"github.com/sirupsen/logrus"

	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/mysql"
)

// Asset for targeting
type Asset struct {
	FID       int32            `db:"id"`
	FHash     string           `db:"hash"`
	FURL      string           `db:"url"`
	FSKU      string           `db:"sku"`
	FBrand    string           `db:"brand"`
	FImg      string           `db:"img"`
	FTitle    string           `db:"title"`
	FPrice    int32            `db:"price"`
	FDiscount int32            `db:"discount"`
	FMime     string           `db:"mime"`
	FCat      mysql.NullString `db:"cat"`
	campaign  entity.Campaign  `db:"-"`
	capping   entity.Capping   `db:"-"`
	assets    []entity.Asset   `db:"-"`
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
	return -1
}

// MimeType of media
func (a *Asset) MimeType() string {
	return a.FMime
}

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
	if a.assets != nil {
		return a.assets
	}
	a.assets = []entity.Asset{
		{
			MimeType: a.FMime,
			Type:     entity.AssetTypeImage,
			SubType:  entity.AssetTypeImageSubTypeMain,
			Width:    250,
			Height:   156,
			Data:     a.FImg,
		},
		{
			MimeType: "text/html",
			Type:     entity.AssetTypeText,
			SubType:  entity.AssetTypeTextSubTypeTitle,
			Len:      int32(utf8.RuneCountInString(a.FTitle)),
			Data:     a.FTitle,
		},
	}
	return a.assets
}

// Cat of product
func (a *Asset) Cat() []openrtb.ContentCategory {
	panic("implement me")
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
func (a *Asset) Price() int32 {
	return a.FPrice
}

// Discount of product if any
func (a *Asset) Discount() int32 {
	return a.FDiscount
}

// ID of item
func (a *Asset) ID() int32 {
	return a.FID
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

// AssetLoader for caching advertiser product
func AssetLoader(_ context.Context) (map[string]kv.Serializable, error) {

	var res []Asset
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
