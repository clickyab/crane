package entities

import (
	"database/sql"
	"fmt"
	"net"
	"time"

	"clickyab.com/crane/crane/entity"
)

// AddImpression insert new impression to daily table
func AddImpression(rh, ref, par, spid, copID string, size, susp int, adid, pubid int64, ip net.IP,
	bid float64, alexa bool, ts time.Time, typ entity.RequestType) error {
	q := fmt.Sprintf(`INSERT INTO impressions%s (
							cp_id,reserved_hash
							w_id,wp_id,app_id,
							ad_id,cop_id,ca_id,
							imp_ipaddress,imp_referaddress,imp_parenturl,
							imp_url,imp_winnerbid,imp_status,
							imp_cookie,imp_alexa,imp_flash,
							imp_time,imp_date,sla_id,slot_id
							) VALUES (
							?,?
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,?
							)`, time.Now().Format("20060102"))
	var err error

	wid := sql.NullInt64{}
	refer := sql.NullString{}
	parent := sql.NullString{}
	if typ == entity.RequestTypeDemand {
		wid = sql.NullInt64{Valid: err == nil, Int64: pubid}
		refer = sql.NullString{Valid: ref != "", String: ref}
		parent = sql.NullString{Valid: par != "", String: par}
	}

	appID := sql.NullInt64{Valid: false}

	// find slot id
	sid, err := FindSlotID(spid, size)
	if err != nil {
		return err
	}
	// insert slot ad
	said, err := InsertSlotAd(sid, adid)
	if err != nil {
		return err
	}
	ca, err := GetAd(adid)
	if err != nil {
		return err
	}
	var alx int
	if alexa {
		alx = 1
	}

	ca.Campaign()
	_, err = NewManager().GetWDbMap().Exec(q,
		ca.Campaign().ID(), rh,
		wid, 0, appID,
		adid, copID, ca.CampaignAdID(),
		ip.String(), refer, parent,
		ca.Target(), bid, susp,
		0, alx, 0,
		ts.Unix(), ts.Format("20060102"), said, sid)
	return err
}
