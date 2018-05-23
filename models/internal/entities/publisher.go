package entities

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/assert"
)

// FindWebsiteByPublicID return publisher id for public id
func FindWebsiteByPublicID(pid int64, supplier entity.Supplier) (entity.Publisher, error) {
	ws := make([]Website, 0)
	_, err := NewManager().GetRDbMap().Select(&ws, `SELECT w_id,
		w_domain,
		w_supplier,
		w_name,
		w_categories,
		w_minbid,
		w_floor_cpm,
		w_fatfinger,
		w_status,
		w_mobad from websites where w_pub_id=?`, pid)
	assert.Nil(err)
	if len(ws) == 0 {
		return nil, errors.New("website not found")
	}
	if supplier.Name() != ws[0].WSupplier {
		return nil, fmt.Errorf("mismatch supplier for domain %s with public id %d. suppliers %s and %s",
			ws[0].WDomain, pid, ws[0].WSupplier, supplier.Name())
	}
	ws[0].Supp = supplier
	return &ws[0], nil
}

// FindAppByAppToken return publisher id for app token
func FindAppByAppToken(token string, supplier entity.Supplier) (entity.Publisher, error) {
	app := make([]App, 0)
	_, err := NewManager().GetRDbMap().Select(&app, `SELECT app_id,app_token,app_name,app_supplier,app_package,app_floor_cpm,app_status,app_cat,app_fatfinger from apps where app_token=?`, token)
	assert.Nil(err)
	if len(app) == 0 {
		return nil, errors.New("app not found")
	}
	if app[0].AppSupplier != supplier.Name() {
		return nil, fmt.Errorf("mismatch supplier for package %s with app token  %s. suppliers %s and %s",
			app[0].AppPackage, token, app[0].AppSupplier, supplier.Name())
	}
	app[0].Supp = supplier
	return &app[0], nil
}

// ErrorNotAllowCreate rise when supplier doesn't allow to add new website or app
var ErrorNotAllowCreate = errors.New("insert not allowed")

// FindOrAddWebsite return publisher id for given supplier,domain
func FindOrAddWebsite(sup entity.Supplier, domain string, pid int64) (entity.Publisher, error) {

	if pid == 0 {
		pid = WebPublicIDGen(sup.Name(), domain)
	}

	w, err := FindWebsiteByPublicID(pid, sup)
	if err == nil {
		return w, nil
	}

	ws := make([]Website, 0)
	_, err = NewManager().GetRDbMap().Select(&ws, `SELECT w_id,
		w_domain,
		w_supplier,
		w_name,
		w_categories,
		w_minbid,
		w_floor_cpm,
		w_fatfinger,
		w_status,
		w_mobad from websites where w_supplier=? and w_domain=?`,
		sup.Name(), domain)
	assert.Nil(err)

	if len(ws) != 0 {
		var tw = ws[0]

		for i := range ws {
			if tw.totalImp() < ws[i].totalImp() {
				tw = ws[i]
			}
		}
		tw.Supp = sup
		return &tw, nil
	}

	if !sup.AllowCreate() {
		return nil, ErrorNotAllowCreate
	}

	tw := Website{
		WSupplier: sup.Name(),
		WDomain:   domain,
		FCTR:      [21]float64{},
		Status:    1,
		Supp:      sup,
		WFloorCpm: sql.NullInt64{Valid: false},
		WMinBid:   sup.DefaultMinBid(),
		WName:     sql.NullString{Valid: true, String: domain},
		CTRStat:   CTRStat{},
	}

	q := `INSERT INTO websites (u_id, w_domain,w_supplier,w_status,created_at,updated_at,w_date, w_pub_id, w_minbid,w_name,w_floor_cpm)
VALUES (?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE
  u_id=VALUES(u_id),w_domain=VALUES(w_domain),w_supplier=VALUES(w_supplier),w_status=VALUES(w_status),
  created_at=VALUES(created_at),updated_at=VALUES(updated_at),w_date=VALUES(w_date), w_id=LAST_INSERT_ID(w_id)`

	t := time.Now()
	res, err := NewManager().GetWDbMap().Exec(q, sup.UserID(),
		tw.WDomain,
		sup.Name(), 1, sql.NullString{Valid: true, String: t.String()}, sql.NullString{Valid: true, String: t.String()},
		int(t.Unix()), WebPublicIDGen(sup.Name(), domain), tw.WMinBid, tw.WName, tw.WFloorCpm)
	if err != nil {
		return nil, err
	}
	tw.WID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &tw, nil

}

// FindOrAddApp return publisher id for given supplier,domain
func FindOrAddApp(sup entity.Supplier, appPackage string, appToken string) (entity.Publisher, error) {

	if appToken == "" {
		appToken = AppPublicIDGen(sup.Name(), appPackage)
	}

	app, err := FindAppByAppToken(appToken, sup)
	if err == nil {
		return app, nil
	}

	apps := make([]App, 0)
	_, err = NewManager().GetRDbMap().Select(&apps, `SELECT app_id,
		app_name,
		app_supplier,
		app_package,
		app_minbid,
		app_status,
		app_floor_cpm,
		app_fatfinger,
		app_cat,
app_token from apps where app_supplier=? and app_package=?`,
		sup.Name(), appPackage)
	assert.Nil(err)

	if len(apps) != 0 {
		var tw = apps[0]

		for i := range apps {
			if tw.totalImp() < apps[i].totalImp() {
				tw = apps[i]
			}
		}
		tw.Supp = sup
		return &tw, nil
	}

	if err != nil && !sup.AllowCreate() {
		return nil, ErrorNotAllowCreate
	}

	q := `INSERT INTO apps (u_id, app_package,app_supplier,app_name,app_status,created_at,updated_at,app_date, app_token)
VALUES (?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE
  u_id=VALUES(u_id),app_package=VALUES(app_package),app_supplier=VALUES(app_supplier),app_name=VALUES(app_name),app_status=VALUES(app_status),
  created_at=VALUES(created_at),updated_at=VALUES(updated_at),app_date=VALUES(app_date), app_id=LAST_INSERT_ID(app_id)`

	t := time.Now()

	tw := App{
		AppSupplier: sup.Name(),
		AppPackage:  appPackage,
		FCTR:        [21]float64{},
		Status:      1,
		Supp:        sup,
		AppFloorCpm: sql.NullInt64{Valid: false},
		AppMinBid:   sup.DefaultMinBid(),
		AppName:     sql.NullString{Valid: true, String: appPackage},
		CTRStat:     CTRStat{},
	}

	res, err := NewManager().GetWDbMap().Exec(q,
		sup.UserID(),
		sql.NullString{Valid: true, String: appPackage},
		sup.Name(),
		sql.NullString{Valid: true, String: appPackage},
		1,
		sql.NullString{Valid: true, String: t.String()},
		sql.NullString{Valid: true, String: t.String()},
		int(t.Unix()), AppPublicIDGen(sup.Name(), appPackage),
	)

	if err != nil {
		return nil, err
	}
	tw.AppID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &tw, nil

}
