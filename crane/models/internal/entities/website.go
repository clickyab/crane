package entities

import (
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/clickyab/services/kv"
)

// website type for website
type Website struct {
	WID         int64          `json:"-" db:"w_id"`
	WDomain     string         `json:"-" db:"w_domain"`
	WSupplier   string         `json:"-" db:"w_supplier"`
	WName       sql.NullString `json:"-" db:"w_name"`
	WCategories SharpArray     `json:"-" db:"w_categories"`
	WMinBid     int64          `json:"-" db:"w_minbid"`
	WFloorCpm   sql.NullInt64  `json:"-" db:"w_floor_cpm"`
	WFatFinger  int            `json:"-" db:"w_fatfinger"`
}

// WebsiteLoader load all confirmed website
func WebsiteLoader(ctx context.Context) (map[string]kv.Serializable, error) {
	var res []Website
	q := `SELECT
			w_id,
			w_domain,
			w_supplier,
			w_name,
			w_categories,
			w_minbid,
			w_floor_cpm,
			w_fatfinger,
		FROM websites WHERE w_status=1`

	_, err := NewManager().GetRDbMap().Select(
		&res,
		q,
	)
	if err != nil {
		return nil, err
	}
	b := make(map[string]kv.Serializable, 0)
	for i := range res {
		b[fmt.Sprintf("%s_%s", res[i].WDomain, res[i].WSupplier)] = &res[i]
	}
	return b, nil
}

func (b *Website) Decode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(b)
}

func (b *Website) Encode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(b)
}
