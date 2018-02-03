package entities

import (
	"context"
	"database/sql"
	"io"

	"github.com/clickyab/services/kv"

	"encoding/gob"

	"fmt"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/simplehash"
)

// App entity
type App struct {
	AppID         int64          `db:"app_id"`
	AppName       sql.NullString `db:"app_name"`
	AppSupplier   string         `db:"app_supplier"`
	AppPackage    string         `db:"app_package"`
	AppMinBid     int64          `db:"app_minbid"`
	Status        int64          `db:"app_status"`
	AppFloorCpm   sql.NullInt64  `db:"app_floor_cpm"`
	AppFatFinger  int            `db:"app_fatfinger"`
	AppCategories SharpArray     `db:"app_cat"`
	AppToken      string         `db:"app_token"`
	Supp          entity.Supplier
	FCTR          [21]float64
	CTRStat

	att map[entity.PublisherAttributes]interface{} `db:"-"`
}

// Attributes return publisher attributes
func (app *App) Attributes() map[entity.PublisherAttributes]interface{} {
	if app.att == nil {
		app.att = make(map[entity.PublisherAttributes]interface{})
		if app.AppFatFinger > 0 {
			app.att[entity.PAFatFinger] = true
		}
	}

	return app.att
}

// Type return type of publisher (app or web)
func (app *App) Type() entity.PublisherType {
	return entity.PublisherTypeApp
}

// Encode for serializable
func (app *App) Encode(iw io.Writer) error {
	return gob.NewEncoder(iw).Encode(app)
}

// Decode for serializable
func (app *App) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(app)
}

// ID return id of app
func (app *App) ID() int64 {
	return app.AppID
}

// FloorCPM return floor cpm of app
func (app *App) FloorCPM() int64 {
	return app.AppFloorCpm.Int64
}

// SoftFloorCPM return soft floor cpm of app
func (app *App) SoftFloorCPM() int64 {
	return app.AppFloorCpm.Int64
}

// Name return name of app
func (app *App) Name() string {
	return app.AppPackage
}

// MinBid return min bid
func (app *App) MinBid() int64 {
	return app.AppMinBid
}

// Supplier return supplier of app
func (app *App) Supplier() entity.Supplier {
	return app.Supp
}

// CTR return ctr of app per size
func (app *App) CTR(size int) float64 {
	if app.FCTR[size] == 0 {
		if app.Supp != nil {
			app.FCTR[size] = app.Supp.DefaultCTR()
		} else {
			app.FCTR[size] = defaultCTR.Float64()
		}
	}
	return app.FCTR[size]
}

// AppLoader load all confirmed apps
func AppLoader(ctx context.Context) (map[string]kv.Serializable, error) {
	b := make(map[string]kv.Serializable)
	return b, nil // Uncomment this line after first time in DEV mode

	const cnt = 10000
	for j := 0; ; j = j + cnt {
		q := fmt.Sprintf(`SELECT app_id,app_token, app_package, app_supplier, app_name, app_cat, app_minbid, app_floor_cpm, app_fatfinger, app_status,
  SUM(imp_1) AS imp1, SUM(imp_2) AS imp2, SUM(imp_3) AS imp3, SUM(imp_4) AS imp4, SUM(imp_5) AS imp5,
  SUM(imp_6) AS imp6, SUM(imp_7) AS imp7, SUM(imp_8) AS imp8, SUM(imp_9) AS imp9, SUM(imp_10) AS imp10,
  SUM(imp_11) AS imp11, SUM(imp_12) AS imp12, SUM(imp_13) AS imp13, SUM(imp_14) AS imp14, SUM(imp_15) AS imp15,
  SUM(imp_16) AS imp16, SUM(imp_17) AS imp17, SUM(imp_18) AS imp18, SUM(imp_19) AS imp19, SUM(imp_20) AS imp20,
  SUM(click_1) AS click1, SUM(click_2) AS click2, SUM(click_3) AS click3, SUM(click_4) AS click4, SUM(click_5) AS click5,
  SUM(click_6) AS click6, SUM(click_7) AS click7, SUM(click_8) AS click8, SUM(click_9) AS click9, SUM(click_10) AS click10,
  SUM(click_11) AS click11, SUM(click_12) AS click12, SUM(click_13) AS click13, SUM(click_14) AS click14, SUM(click_15) AS click15,
  SUM(click_16) AS click16, SUM(click_17) AS click17, SUM(click_18) AS click18, SUM(click_19) AS click19, SUM(click_20) AS click20
  FROM apps
  LEFT JOIN ctr_stat ON app_id=pub_id AND pub_type=? AND date BETWEEN DATE_SUB(NOW(), INTERVAL 2 DAY) AND NOW()
  GROUP BY app_id LIMIT %d, %d`, j, j+cnt)

		var res []App
		if _, err := NewManager().GetRDbMap().Select(&res, q, "app"); err != nil {
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
			k := fmt.Sprintf("%s/%s", res[i].AppSupplier, res[i].AppPackage)
			if d, ok := b[k]; ok {
				if n.totalImp() < d.(*App).totalImp() {
					n = d.(*App)
				}
			}
			b[k] = n
		}
	}
	return b, nil
}

func (app *App) totalImp() (res int64) {
	if app.Status != 1 {
		return -1
	}
	return app.Impression1.Int64 +
		app.Impression2.Int64 +
		app.Impression3.Int64 +
		app.Impression4.Int64 +
		app.Impression5.Int64 +
		app.Impression6.Int64 +
		app.Impression7.Int64 +
		app.Impression8.Int64 +
		app.Impression9.Int64 +
		app.Impression10.Int64 +
		app.Impression11.Int64 +
		app.Impression12.Int64 +
		app.Impression13.Int64 +
		app.Impression14.Int64 +
		app.Impression15.Int64 +
		app.Impression16.Int64 +
		app.Impression17.Int64 +
		app.Impression18.Int64 +
		app.Impression19.Int64 +
		app.Impression20.Int64
}

// AppPublicIDGen generate app token from supplier name and package
func AppPublicIDGen(sup, appPackage string) string {
	return simplehash.MD5(sup + "/" + appPackage)
}
