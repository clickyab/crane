package entities

import (
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"io"

	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/kv"
)

// website type for website
type Website struct {
	WID         int64          `db:"w_id"`
	WDomain     string         `db:"w_domain"`
	WSupplier   string         `db:"w_supplier"`
	WName       sql.NullString `db:"w_name"`
	WCategories SharpArray     `db:"w_categories"`
	WMinBid     int64          `db:"w_minbid"`
	WFloorCpm   sql.NullInt64  `db:"w_floor_cpm"`
	CTRStat

	FCTR [21]float64
}

func (w *Website) CTR(size int) float64 {
	return w.FCTR[size]
}

func (w *Website) ID() int64 {
	return w.WID
}

// soft floor and floor cpm are the same
func (w *Website) FloorCPM() int64 {
	return w.WFloorCpm.Int64
}

func (w *Website) SoftFloorCPM() int64 {
	return w.WFloorCpm.Int64
}

func (w *Website) Name() string {
	return w.WName.String
}

func (w *Website) Attributes() map[string]interface{} {
	return map[string]interface{}{}
}

func (w *Website) BIDType() entity.BIDType {
	return entity.BIDTypeCPC
}

func (w *Website) MinBid() int64 {
	return w.WMinBid
}

// just supporting banner for now
func (w *Website) AcceptedTypes() []entity.AdType {
	return []entity.AdType{entity.AdTypeBanner}
}

func (w *Website) Supplier() string {
	return w.WSupplier
}

// WebsiteLoader load all confirmed website
func WebsiteLoader(ctx context.Context) (map[string]kv.Serializable, error) {
	q := `SELECT w_id, w_domain, w_supplier, w_name, w_categories, w_minbid, w_floor_cpm, w_fatfinger,
  SUM(imp_1) AS imp1, SUM(imp_2) AS imp2, SUM(imp_3) AS imp3, SUM(imp_4) AS imp4, SUM(imp_5) AS imp5,
  SUM(imp_6) AS imp6, SUM(imp_7) AS imp7, SUM(imp_8) AS imp8, SUM(imp_9) AS imp9, SUM(imp_10) AS imp10,
  SUM(imp_11) AS imp11, SUM(imp_12) AS imp12, SUM(imp_13) AS imp13, SUM(imp_14) AS imp14, SUM(imp_15) AS imp15,
  SUM(imp_16) AS imp16, SUM(imp_17) AS imp17, SUM(imp_18) AS imp18, SUM(imp_19) AS imp19, SUM(imp_20) AS imp20,
  SUM(click_1) AS click1, SUM(click_2) AS click2, SUM(click_3) AS click3, SUM(click_4) AS click4, SUM(click_5) AS click5,
  SUM(click_6) AS click6, SUM(click_7) AS click7, SUM(click_8) AS click8, SUM(click_9) AS click9, SUM(click_10) AS click10,
  SUM(click_11) AS click11, SUM(click_12) AS click12, SUM(click_13) AS click13, SUM(click_14) AS click14, SUM(click_15) AS click15,
  SUM(click_16) AS click16, SUM(click_17) AS click17, SUM(click_18) AS click18, SUM(click_19) AS click19, SUM(click_20) AS click20
  FROM websites
  INNER JOIN ctr_stat ON w_id=pub_id
  WHERE date BETWEEN DATE_SUB(NOW(), INTERVAL 2 DAY) AND NOW()
  AND w_status=1 AND pub_type=?
  GROUP BY w_id`

	var res []Website
	if _, err := NewManager().GetRDbMap().Select(&res, q, "web"); err != nil {
		return nil, err
	}

	b := make(map[string]kv.Serializable, 0)
	for i := range res {
		res[i].FCTR = [21]float64{}
		res[i].FCTR[1] = res[i].Click1 / res[i].Impression1 * 100
		res[i].FCTR[2] = res[i].Click2 / res[i].Impression2 * 100
		res[i].FCTR[3] = res[i].Click3 / res[i].Impression3 * 100
		res[i].FCTR[4] = res[i].Click4 / res[i].Impression4 * 100
		res[i].FCTR[5] = res[i].Click5 / res[i].Impression5 * 100
		res[i].FCTR[6] = res[i].Click6 / res[i].Impression6 * 100
		res[i].FCTR[7] = res[i].Click7 / res[i].Impression7 * 100
		res[i].FCTR[8] = res[i].Click8 / res[i].Impression8 * 100
		res[i].FCTR[9] = res[i].Click9 / res[i].Impression9 * 100
		res[i].FCTR[10] = res[i].Click10 / res[i].Impression10 * 100
		res[i].FCTR[11] = res[i].Click11 / res[i].Impression11 * 100
		res[i].FCTR[12] = res[i].Click12 / res[i].Impression12 * 100
		res[i].FCTR[13] = res[i].Click13 / res[i].Impression13 * 100
		res[i].FCTR[14] = res[i].Click14 / res[i].Impression14 * 100
		res[i].FCTR[15] = res[i].Click15 / res[i].Impression15 * 100
		res[i].FCTR[16] = res[i].Click16 / res[i].Impression16 * 100
		res[i].FCTR[17] = res[i].Click17 / res[i].Impression17 * 100
		res[i].FCTR[18] = res[i].Click18 / res[i].Impression18 * 100
		res[i].FCTR[19] = res[i].Click19 / res[i].Impression19 * 100
		res[i].FCTR[20] = res[i].Click20 / res[i].Impression20 * 100

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
