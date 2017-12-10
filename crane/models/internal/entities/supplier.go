package entities

import (
	"context"
	"time"

	"io"

	"encoding/gob"

	"database/sql"

	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/kv"
)

// Supplier is the supplier structure
type Supplier struct {
	ID               int64         `db:"id"`
	FName            string        `db:"name"`
	FToken           string        `db:"token"`
	UserID           sql.NullInt64 `db:"user_id"`
	DefaultFloor     int64         `db:"default_floor"`
	DefaultSoftFloor int64         `db:"default_soft_floor"`
	DefaultMinBID    int64         `db:"default_min_bid"`
	FBIDType         string        `db:"bid_type"`
	FDefaultCTR      float64       `db:"default_ctr"`
	Tiny             int           `db:"tiny_mark"`
	CreatedAt        time.Time     `db:"created_at"`
	UpdatedAt        time.Time     `db:"updated_at"`
}

// TinyMark show tiny clickyb mark
func (s *Supplier) TinyMark() bool {
	return s.Tiny > 0
}

// AllowCreate allow create new site on demand?
func (s *Supplier) AllowCreate() bool {
	return s.UserID.Valid
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

// SupplierLoader load all confirmed website
func SupplierLoader(ctx context.Context) (map[string]kv.Serializable, error) {
	q := `SELECT * FROM suppliers`

	var res []Supplier
	if _, err := NewManager().GetRDbMap().Select(&res, q, "web"); err != nil {
		return nil, err
	}

	b := make(map[string]kv.Serializable)
	for i := range res {
		b[res[i].FToken] = &res[i]
	}
	return b, nil
}
