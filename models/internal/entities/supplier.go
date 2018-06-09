package entities

import (
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"io"
	"strings"
	"time"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/mysql"
)

// Supplier is the supplier structure
type Supplier struct {
	FName         string                 `db:"name"`
	FToken        string                 `db:"token"`
	FUserID       sql.NullInt64          `db:"user_id"`
	FSoftFloorCPM mysql.GenericJSONField `db:"soft_floor_cpm"`
	FSoftFloorCPC mysql.GenericJSONField `db:"soft_floor_cpc"`
	DefaultMinBID int64                  `db:"default_min_bid"`
	FBIDType      string                 `db:"bid_type"`
	FDefaultCTR   mysql.GenericJSONField `db:"ctr"`
	Tiny          int                    `db:"tiny_mark"`
	FShowDomain   string                 `db:"show_domain"`
	CreatedAt     time.Time              `db:"created_at"`
	UpdatedAt     time.Time              `db:"updated_at"`
	FRate         int                    `db:"rate"`
	FTinyLogo     string                 `db:"tiny_logo"`
	FTinyURL      string                 `db:"tiny_url"`
	FShare        int                    `db:"share"`
	FMarkup       bool                   `db:"markup"`
	strategy      entity.Strategy        `db:"-"`
}

// SoftFloorCPC based on pub type and request type
func (s Supplier) SoftFloorCPC(adType, pubType string) int64 {
	key := fmt.Sprintf("%s_%s", pubType, adType)
	if val, ok := s.FSoftFloorCPC[key]; ok {
		if x, ok := val.(float64); ok {
			return int64(x)
		}
	}
	panic("[BUG]supplier not support proper floor cpm")
}

// Strategy of supplier can be cpm, cpc or both
func (s Supplier) Strategy() entity.Strategy {
	if s.strategy != 0 {
		return s.strategy
	}
	s.strategy = entity.GetStrategy(strings.Split(s.FBIDType, ","))
	if s.strategy == 0 {
		// fall back to default strategy in case of data error
		s.strategy = entity.StrategyCPM
	}
	return s.strategy
}

// SoftFloorCPM based on pub type and request type
func (s Supplier) SoftFloorCPM(adType string, pubType string) int64 {
	key := fmt.Sprintf("%s_%s", pubType, adType)
	if val, ok := s.FSoftFloorCPM[key]; ok {
		if x, ok := val.(float64); ok {
			return int64(x)
		}
	}
	panic("[BUG]supplier not support proper floor cpm")
}

// DefaultCTR get ctr by ad type and publisher type
func (s Supplier) DefaultCTR(adType string, pubType string) float64 {
	key := fmt.Sprintf("%s_%s", pubType, adType)
	if val, ok := s.FDefaultCTR[key]; ok {
		if x, ok := val.(float64); ok {
			return x
		}
	}
	panic("[BUG]supplier not support proper default ctr")
}

// TinyLogo will be the url to the logo (ex: //clickyab.com/tiny.png)
func (s Supplier) TinyLogo() string {
	return s.FTinyLogo
}

// TinyURL is the link of ancher tag of tiny (ex: http://clickyab.com/?ref=tiny)
func (s Supplier) TinyURL() string {
	return s.FTinyURL
}

// Rate return ratio currency conversion to IRR
func (s Supplier) Rate() int {
	return s.FRate
}

// ShowDomain is a domain that all links are generated against it
func (s Supplier) ShowDomain() string {
	return s.FShowDomain
}

// UserID return user id of supplier
func (s *Supplier) UserID() int64 {
	return s.FUserID.Int64
}

// TinyMark show tiny clickyab mark
func (s *Supplier) TinyMark() bool {
	return s.Tiny > 0
}

// AllowCreate allow create new site on demand?
func (s *Supplier) AllowCreate() bool {
	return s.FUserID.Valid
}

// Encode is the encode function for serialize object in io writer
func (s *Supplier) Encode(iw io.Writer) error {
	return gob.NewEncoder(iw).Encode(s)
}

// Decode try to decode object from io reader
func (s *Supplier) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(s)
}

// Name of supplier
func (s *Supplier) Name() string {
	return s.FName
}

// Token is used for finding supplier
func (s *Supplier) Token() string {
	return s.FToken
}

// Markup return markup status
func (s *Supplier) Markup() bool {
	return s.FMarkup
}

// DefaultMinBid is the default min bid
func (s *Supplier) DefaultMinBid() int64 {
	return s.DefaultMinBID
}

// Share of this supplier
func (s *Supplier) Share() int {
	return s.FShare
}

var (
	supQuery = `SELECT markup,name,token,user_id,soft_floor_cpm, soft_floor_cpc, default_min_bid,bid_type,ctr,tiny_mark,
show_domain,created_at,updated_at,rate,tiny_logo,tiny_url,share FROM suppliers`
)

// SupplierLoader load all confirmed website
func SupplierLoader(_ context.Context) (map[string]kv.Serializable, error) {

	var res []Supplier
	if _, err := NewManager().GetRDbMap().Select(&res, supQuery); err != nil {
		return nil, err
	}

	b := make(map[string]kv.Serializable)
	for i := range res {
		b[res[i].FToken] = &res[i]
	}
	return b, nil
}

// SupplierLoaderByName load all confirmed website
func SupplierLoaderByName(_ context.Context) (map[string]kv.Serializable, error) {

	var res []Supplier
	if _, err := NewManager().GetRDbMap().Select(&res, supQuery); err != nil {
		return nil, err
	}

	b := make(map[string]kv.Serializable)
	for i := range res {
		b[res[i].FName] = &res[i]
	}
	return b, nil
}
