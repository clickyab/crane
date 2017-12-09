package entities

import (
	"database/sql"
	"io"

	"fmt"

	"encoding/gob"

	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/kv"
)

const appsDBName = `apps`

// App is the applications structure
type App struct {
	FID         int64         `db:"app_id"`
	UserID      int64         `db:"u_id"`
	Token       string        `db:"app_token"`
	FName       string        `db:"app_name"`
	EnName      string        `db:"en_app_name"`
	Package     string        `db:"app_package"`
	FSupplier   string        `db:"app_supplier"`
	FFloorCPM   sql.NullInt64 `db:"app_floor_cpm"`
	Status      int           `db:"app_status"`
	TodayCTR    int64         `db:"app_today_ctr"`
	TodayIMPs   int64         `db:"app_today_imps"`
	TodayClicks int64         `db:"app_today_clicks"`
	AppMinBid   int64         `db:"app_minbid"`
	Category    SharpArray    `db:"app_cat"`
	CTRStat

	FCTR [21]float64
}

// CTR return the ctr of this app based on the size
func (a App) CTR(size int) float64 {
	return a.FCTR[size]
}

// Decode the decoder
func (a App) Decode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(a)
}

// Encode encoder
func (a App) Encode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(a)
}

// ID app id
func (a *App) ID() int64 {
	return a.FID
}

// FloorCPM needs to be zero if its not set in db
func (a *App) FloorCPM() int64 {
	return a.FFloorCPM.Int64
}

// SoftFloorCPM soft returns actual floor
func (a *App) SoftFloorCPM() int64 {
	return a.FFloorCPM.Int64
}

// Name app name
func (a *App) Name() string {
	return a.FName
}

// BIDType the bid type of this app
func (a *App) BIDType() entity.BIDType {
	return entity.BIDTypeCPC
}

// MinBid is the minimum bid
func (a *App) MinBid() int64 {
	return a.AppMinBid
}

// Supplier is the supplier object
func (a *App) Supplier() string {
	return a.FSupplier
}

// AppLoader is the loader for cache
func AppLoader() (map[string]kv.Serializable, error) {
	q := fmt.Sprintf(`SELECT
  app_id, u_id, app_token, app_name, en_app_name, app_package, app_supplier, app_floor_cpm,
  app_status, app_today_ctr, app_today_imps, app_today_clicks, app_cat,app_minbid,
  SUM(imp_1) AS imp1, SUM(imp_2) AS imp2, SUM(imp_3) AS imp3, SUM(imp_4) AS imp4, SUM(imp_5) AS imp5,
  SUM(imp_6) AS imp6, SUM(imp_7) AS imp7, SUM(imp_8) AS imp8, SUM(imp_9) AS imp9, SUM(imp_10) AS imp10,
  SUM(imp_11) AS imp11, SUM(imp_12) AS imp12, SUM(imp_13) AS imp13, SUM(imp_14) AS imp14, SUM(imp_15) AS imp15,
  SUM(imp_16) AS imp16, SUM(imp_17) AS imp17, SUM(imp_18) AS imp18, SUM(imp_19) AS imp19, SUM(imp_20) AS imp20,
  SUM(click_1) AS click1, SUM(click_2) AS click2, SUM(click_3) AS click3, SUM(click_4) AS click4, SUM(click_5) AS click5,
  SUM(click_6) AS click6, SUM(click_7) AS click7, SUM(click_8) AS click8, SUM(click_9) AS click9, SUM(click_10) AS click10,
  SUM(click_11) AS click11, SUM(click_12) AS click12, SUM(click_13) AS click13, SUM(click_14) AS click14, SUM(click_15) AS click15,
  SUM(click_16) AS click16, SUM(click_17) AS click17, SUM(click_18) AS click18, SUM(click_19) AS click19, SUM(click_20) AS click20
  FROM %s INNER JOIN ctr_stat ON app_id=pub_id WHERE
  date BETWEEN DATE_SUB(NOW(), INTERVAL 2 DAY) AND NOW()
  AND app_status=1 AND pub_type=?
  GROUP BY app_id`, appsDBName)
	var res []App
	_, err := NewManager().GetRDbMap().Select(&res, q, "app")
	if err != nil {
		return nil, err
	}

	resp := map[string]kv.Serializable{}
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

		resp[res[i].FName] = res[i]
	}

	return resp, nil
}
