package entities

import (
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"io"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/mysql"
	"github.com/clickyab/services/simplehash"
)

var (
	supplierName = config.RegisterString("crane.supplier.clickyab.name", "clickyab", "name of our supplier name")
)

// Website type for website
type Website struct {
	WID         int64                  `db:"w_id"`
	WDomain     string                 `db:"w_domain"`
	WSupplier   string                 `db:"w_supplier"`
	WName       sql.NullString         `db:"w_name"`
	WCategories SharpArray             `db:"w_categories"`
	WMinBid     int64                  `db:"w_minbid"`
	WFloorCpm   sql.NullInt64          `db:"w_floor_cpm"`
	WFatFinger  int                    `db:"w_fatfinger"`
	Status      int                    `db:"w_status"`
	MobAd       int                    `db:"w_mobad"`
	PublicID    int64                  `db:"w_pub_id"`
	WMinCPC     mysql.GenericJSONField `db:"w_min_cpc"`
	CTRStat
	Supp entity.Supplier `db:"-"`
	FCTR [21]float64

	att map[entity.PublisherAttributes]interface{}
}

// Attributes return publisher attributes
func (w *Website) Attributes() map[entity.PublisherAttributes]interface{} {
	if w.att == nil {
		w.att = make(map[entity.PublisherAttributes]interface{})
		if w.MobAd > 0 {
			w.att[entity.PAMobileAd] = true
		}
		if w.WFatFinger > 0 {
			w.att[entity.PAFatFinger] = true
		}
	}
	return w.att
}

// Type return type of publisher (app or web)
func (w *Website) Type() entity.PublisherType {
	return entity.PublisherTypeWeb
}

// CTR return the ctr based on size of this website
func (w *Website) CTR(size int) float64 {
	if w.FCTR[size] == 0 {
		w.FCTR[size] = -1
	}
	return w.FCTR[size]
}

func (w *Website) totalImp() (res int64) {
	if w.Status != 1 {
		return -1
	}
	return w.Impression1.Int64 +
		w.Impression2.Int64 +
		w.Impression3.Int64 +
		w.Impression4.Int64 +
		w.Impression5.Int64 +
		w.Impression6.Int64 +
		w.Impression7.Int64 +
		w.Impression8.Int64 +
		w.Impression9.Int64 +
		w.Impression10.Int64 +
		w.Impression11.Int64 +
		w.Impression12.Int64 +
		w.Impression13.Int64 +
		w.Impression14.Int64 +
		w.Impression15.Int64 +
		w.Impression16.Int64 +
		w.Impression17.Int64 +
		w.Impression18.Int64 +
		w.Impression19.Int64 +
		w.Impression20.Int64
}

// ID return the website id
func (w *Website) ID() int64 {
	return w.WID
}

// FloorCPM soft floor and floor cpm are the same
func (w *Website) FloorCPM() int64 {
	return w.WFloorCpm.Int64
}

// SoftFloorCPM the soft flor ans we need to accept this first then fall back to floorcpm
func (w *Website) SoftFloorCPM() int64 {
	return w.WFloorCpm.Int64
}

// Name of the website
func (w *Website) Name() string {
	return w.WDomain
}

// MinBid return the minimum bid accepted for this
func (w *Website) MinBid() int64 {
	return w.WMinBid
}

// Supplier of this website
func (w *Website) Supplier() entity.Supplier {
	return w.Supp
}

// Categories return publisher categories
func (w *Website) Categories() []string {
	var res = make([]string, 0)
	for i := range w.WCategories {
		res = append(res, "IAB"+w.WCategories[i])
	}
	return res
}

// WebsiteLoaderGen load all confirmed website
func WebsiteLoaderGen(name bool) func(ctx context.Context) (map[string]kv.Serializable, error) {
	return func(ctx context.Context) (map[string]kv.Serializable, error) {
		b := make(map[string]kv.Serializable)
		return b, nil // Uncomment this line after first time in DEV mode

		where := ""
		if !name {
			where = fmt.Sprintf(" WHERE w_supplier='%s' ", supplierName)
		}
		const cnt = 10000
		for j := 0; ; j = j + cnt {
			q := fmt.Sprintf(`SELECT w_id, w_domain, w_supplier, w_name, w_categories, w_minbid, w_floor_cpm, w_fatfinger, w_status,w_mobad,w_pub_id,w_min_cpc,
  SUM(imp_1) AS imp1, SUM(imp_2) AS imp2, SUM(imp_3) AS imp3, SUM(imp_4) AS imp4, SUM(imp_5) AS imp5,
  SUM(imp_6) AS imp6, SUM(imp_7) AS imp7, SUM(imp_8) AS imp8, SUM(imp_9) AS imp9, SUM(imp_10) AS imp10,
  SUM(imp_11) AS imp11, SUM(imp_12) AS imp12, SUM(imp_13) AS imp13, SUM(imp_14) AS imp14, SUM(imp_15) AS imp15,
  SUM(imp_16) AS imp16, SUM(imp_17) AS imp17, SUM(imp_18) AS imp18, SUM(imp_19) AS imp19, SUM(imp_20) AS imp20,
  SUM(click_1) AS click1, SUM(click_2) AS click2, SUM(click_3) AS click3, SUM(click_4) AS click4, SUM(click_5) AS click5,
  SUM(click_6) AS click6, SUM(click_7) AS click7, SUM(click_8) AS click8, SUM(click_9) AS click9, SUM(click_10) AS click10,
  SUM(click_11) AS click11, SUM(click_12) AS click12, SUM(click_13) AS click13, SUM(click_14) AS click14, SUM(click_15) AS click15,
  SUM(click_16) AS click16, SUM(click_17) AS click17, SUM(click_18) AS click18, SUM(click_19) AS click19, SUM(click_20) AS click20
  FROM websites
  LEFT JOIN ctr_stat ON w_id=pub_id AND pub_type=? AND date BETWEEN DATE_SUB(NOW(), INTERVAL 2 DAY) AND NOW() %s
  GROUP BY w_id LIMIT %d, %d`, where, j, j+cnt)

			var res []Website
			if _, err := NewManager().GetRDbMap().Select(&res, q, "web"); err != nil {
				return nil, err
			}
			if len(res) == 0 {
				break
			}

			for i := range res {
				res[i].FCTR = [21]float64{}
				res[i].FCTR[1] = calc(res[i].Impression1, res[i].Click1)
				res[i].FCTR[2] = calc(res[i].Impression2, res[i].Click2)
				res[i].FCTR[3] = calc(res[i].Impression3, res[i].Click3)
				res[i].FCTR[4] = calc(res[i].Impression4, res[i].Click4)
				res[i].FCTR[5] = calc(res[i].Impression5, res[i].Click5)
				res[i].FCTR[6] = calc(res[i].Impression6, res[i].Click6)
				res[i].FCTR[7] = calc(res[i].Impression7, res[i].Click7)
				res[i].FCTR[8] = calc(res[i].Impression8, res[i].Click8)
				res[i].FCTR[9] = calc(res[i].Impression9, res[i].Click9)
				res[i].FCTR[10] = calc(res[i].Impression10, res[i].Click10)
				res[i].FCTR[11] = calc(res[i].Impression11, res[i].Click11)
				res[i].FCTR[12] = calc(res[i].Impression12, res[i].Click12)
				res[i].FCTR[13] = calc(res[i].Impression13, res[i].Click13)
				res[i].FCTR[14] = calc(res[i].Impression14, res[i].Click14)
				res[i].FCTR[15] = calc(res[i].Impression15, res[i].Click15)
				res[i].FCTR[16] = calc(res[i].Impression16, res[i].Click16)
				res[i].FCTR[17] = calc(res[i].Impression17, res[i].Click17)
				res[i].FCTR[18] = calc(res[i].Impression18, res[i].Click18)
				res[i].FCTR[19] = calc(res[i].Impression19, res[i].Click19)
				res[i].FCTR[20] = calc(res[i].Impression20, res[i].Click20)
				n := &res[i]
				k := fmt.Sprintf("%s/%s", res[i].WSupplier, res[i].WDomain)
				if !name {
					k = fmt.Sprint(res[i].PublicID)
				}
				if d, ok := b[k]; ok {
					if n.totalImp() < d.(*Website).totalImp() {
						n = d.(*Website)
					}
				}
				b[k] = n
			}
		}
		return b, nil
	}
}

// Encode is the encode function for serialize object in io writer
func (w *Website) Encode(iw io.Writer) error {
	return gob.NewEncoder(iw).Encode(w)
}

// Decode try to decode object from io reader
func (w *Website) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(w)
}

// MinCPC return min cpc for this site
func (w *Website) MinCPC(adType string) float64 {
	if val, ok := w.WMinCPC[adType]; ok {
		if x, ok := val.(float64); ok {
			return x
		}
	}
	return 0
}

// WebPublicIDGen generate public id from supplier name and domain
func WebPublicIDGen(sup, domain string) int64 {
	res := simplehash.CRC64(sup + "/" + domain)
	if res > 0 {
		res *= -1
	}
	return res
}
