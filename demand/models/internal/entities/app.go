package entities

import (
	"database/sql"

	"clickyab.com/crane/demand/entity"
)

// App entity
type App struct {
	AppID         int64          `db:"app_id"`
	AppName       sql.NullString `db:"app_name"`
	AppSupplier   int64          `db:"app_supplier"`
	AppPackage    int64          `db:"app_package"`
	AppMinBid     int64          `db:"app_minbid"`
	Status        int64          `db:"app_status"`
	AppFloorCpm   sql.NullInt64  `db:"app_floor_cpm"`
	AppFatFinger  int            `db:"app_fatfinger"`
	AppCategories SharpArray     `db:"app_cat"`

	Supp entity.Supplier
	FCTR [21]float64
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
	return app.AppName.String
}

// BIDType return bid type cpc,cpm
func (app *App) BIDType() entity.BIDType {
	return entity.BIDTypeCPC
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
