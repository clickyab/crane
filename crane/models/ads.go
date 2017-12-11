package models

import (
	"database/sql"
	"errors"
	"fmt"
	"hash/crc64"
	"net"
	"time"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/models/internal/entities"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/pool"
)

var userID = config.RegisterInt64("crane.model.db.user_id", 0, "default user id for workers")

var ads pool.Interface

// GetAds return all ads in system
func GetAds() []entity.Advertise {
	data := ads.All()
	all := make([]entity.Advertise, len(data))
	var c int
	for i := range data {
		all[c] = data[i].(entity.Advertise)
		c++
	}

	return all
}

// GetAd try to get advertise based on its id
func GetAd(adID int64) (entity.Advertise, error) {
	ad, err := ads.Get(fmt.Sprint(adID), &entities.Advertise{})
	if err != nil {
		x, err := entities.GetAd(adID)
		if err != nil {
			return nil, err
		}
		return x, nil
	}
	return ad.(entity.Advertise), nil
}

var ErrorNotAllowCreate = errors.New("insert not allowed")

// FindPublisherId return publisher id for given supplier,domain
func FindPublisherId(sup, domain string) (int64, error) {
	osup, err := GetSupplierByName(sup)
	if err != nil {
		return 0, err
	}
	osup.AllowCreate()
	p, err := GetWebSite(osup, domain)
	if err == nil {
		return p.ID(), nil
	}
	if !osup.AllowCreate() {
		return 0, ErrorNotAllowCreate
	}
	crc := crc64.New(crc64.MakeTable(crc64.ECMA))
	crc.Write([]byte(sup + "/" + domain))
	sum := int64(crc.Sum64())

	q := `INSERT INTO websites (u_id, w_domain,w_supplier,w_status,created_at,updated_at,w_date, w_pub_id )
VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE
  u_id=VALUES(u_id),w_domain=VALUES(w_domain),w_supplier=VALUES(w_supplier),w_status=VALUES(w_status),
  created_at=VALUES(created_at),updated_at=VALUES(updated_at),w_date=VALUES(w_date), w_id=LAST_INSERT_ID(w_id)
	`
	t := time.Now()
	res, err := entities.NewManager().GetWDbMap().Exec(q, userID.Int64(),
		sql.NullString{Valid: true, String: domain},
		sup, 1, sql.NullString{Valid: true, String: t.String()}, sql.NullString{Valid: true, String: t.String()},
		int(t.Unix()), sum,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// FindSlotID return slot id for given public-id, slot-size
func FindSlotID(pid string, s int) (int64, error) {
	q := `INSERT INTO slots (slot_public_id, slot_size) VALUES (?, ?) ON DUPLICATE KEY UPDATE
slot_public_id=VALUES(slot_public_id), slot_size=VALUES(slot_size), slot_id=LAST_INSERT_ID(slot_id)`
	m, err := entities.NewManager().GetWDbMap().Exec(q, pid, s)
	if err != nil {
		return 0, err
	}
	return m.LastInsertId()
}

// InsertSlotAd insert into slot to slots_ads table
func InsertSlotAd(sid, adid int64) (int64, error) {
	query := `INSERT INTO slots_ads (slot_id, ad_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE
sla_id=LAST_INSERT_ID(sla_id)`
	res, err := entities.NewManager().GetWDbMap().Exec(
		query,
		sid,
		adid,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func AddImpression(sup, pub, ref, par, pid, copID string, size, susp int, adid int64, ip net.IP,
	bid float64, alexa bool, ts time.Time, typ entity.RequestType) error {
	tw, err := FindPublisherId(sup, pub)
	wid := sql.NullInt64{}
	refer := sql.NullString{}
	parent := sql.NullString{}
	if typ == entity.RequestTypeDemand {
		wid = sql.NullInt64{Valid: err == nil, Int64: tw}
		refer = sql.NullString{Valid: ref != "", String: ref}
		parent = sql.NullString{Valid: par != "", String: par}
	}
	appID := sql.NullInt64{Valid: false}
	q := fmt.Sprintf(`INSERT INTO impressions%s (
							cp_id,
							w_id,wp_id,app_id,
							ad_id,cop_id,ca_id,
							imp_ipaddress,imp_referaddress,imp_parenturl,
							imp_url,imp_winnerbid,imp_status,
							imp_cookie,imp_alexa,imp_flash,
							imp_time,imp_date,sla_id,slot_id
							) VALUES (
							?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,?
							)`, time.Now().Format("20060102"))

	// find slot id
	sid, err := FindSlotID(pid, size)
	if err != nil {
		return err
	}
	// insert slot ad
	said, err := InsertSlotAd(sid, adid)
	if err != nil {
		return err
	}
	_, err = entities.NewManager().GetWDbMap().Exec(q,
		"",
		wid, 0, appID,
		adid, copID, 0,
		ip.String(), refer, parent,
		"", bid, susp,
		0, alexa, 0,
		ts.Unix(), ts.Format("20060102"),
		said, sid)
	return err
}
