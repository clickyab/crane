package entities

import (
	"context"
	"time"

	"io"

	"encoding/gob"

	"database/sql"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/kv"
)

// Supplier is the supplier structure
type Supplier struct {
	FName            string        `db:"name"`
	FToken           string        `db:"token"`
	FUserID          sql.NullInt64 `db:"user_id"`
	DefaultFloor     int64         `db:"default_floor"`
	DefaultSoftFloor int64         `db:"default_soft_floor"`
	DefaultMinBID    int64         `db:"default_min_bid"`
	FBIDType         string        `db:"bid_type"`
	FDefaultCTR      float64       `db:"default_ctr"`
	Tiny             int           `db:"tiny_mark"`
	FShowDomain      string        `db:"show_domain"`
	CreatedAt        time.Time     `db:"created_at"`
	UpdatedAt        time.Time     `db:"updated_at"`
	FRate            int           `db:"rate"`
	FTinyLogo        string        `db:"tiny_logo"`
	FTinyURL         string        `db:"tiny_url"`
	FUnderFloor      int           `db:"under_floor"`
	FShare           int           `db:"share"`
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

// DefaultFloorCPM is the default floor for new sites
func (s *Supplier) DefaultFloorCPM() int64 {
	return s.DefaultFloor
}

// DefaultSoftFloorCPM is the default floor for new sites
func (s *Supplier) DefaultSoftFloorCPM() int64 {
	return s.DefaultSoftFloor
}

// DefaultMinBid is the default min bid
func (s *Supplier) DefaultMinBid() int64 {
	return s.DefaultMinBID
}

// BidType return this supplier bid type
func (s *Supplier) BidType() entity.BIDType {
	switch s.FBIDType {
	case "cpm":
		return entity.BIDTypeCPM
	default:
		return entity.BIDTypeCPC
	}
}

// DefaultCTR used for this website no data ctr calculation
func (s *Supplier) DefaultCTR() float64 {
	return s.FDefaultCTR
}

// UnderFloor means that this supplier allow to pass underfloor value.
// normally used only for clickyab
func (s *Supplier) UnderFloor() bool {
	return s.FUnderFloor != 0
}

// Share of this supplier
func (s *Supplier) Share() int {
	return s.FShare
}

var (
	supQuery = `SELECT name,token,user_id,default_floor,default_soft_floor,default_min_bid,bid_type,default_ctr,tiny_mark,
show_domain,created_at,updated_at,rate,tiny_logo,tiny_url,under_floor,share FROM suppliers`
)

// SupplierLoader load all confirmed website
func SupplierLoader(ctx context.Context) (map[string]kv.Serializable, error) {

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
func SupplierLoaderByName(ctx context.Context) (map[string]kv.Serializable, error) {

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
