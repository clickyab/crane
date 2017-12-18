package entities

import (
	"database/sql"
	"fmt"
	"net"
	"time"

	"clickyab.com/crane/demand/entity"
)

// AddImpression insert new impression to daily table
// TODO : multiple insert per one query
func AddImpression(rh, ref, par, spid, copID string, size, susp int, adid, pubid int64, ip net.IP,
	bid float64, alexa bool, ts time.Time, typ entity.RequestType) error {
	var err error

	wID := sql.NullInt64{}
	appID := sql.NullInt64{}
	refer := sql.NullString{}
	parent := sql.NullString{}
	if typ == entity.RequestTypeDemand {
		// TODO : check for publisher type in demand too
		wID = sql.NullInt64{Valid: err == nil, Int64: pubid}
		refer = sql.NullString{Valid: ref != "", String: ref}
		parent = sql.NullString{Valid: par != "", String: par}
	}

	var sID int64

	// find slot id
	if wID.Valid {
		sID, err = FindWebSlotID(spid, wID.Int64, size)
	} else if appID.Valid {
		sID, err = FindAppSlotID(spid, appID.Int64, size)
	}
	if err != nil {
		return err
	}

	// find slot ad
	said, err := FindSlotAd(sID, adid)
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
	q := fmt.Sprintf(`INSERT INTO impressions%s (
							cp_id,reserved_hash,ad_size,
							w_id,wp_id,app_id,
							ad_id,cop_id,ca_id,
							imp_ipaddress,imp_referaddress,imp_parenturl,
							imp_url,imp_winnerbid,imp_status,
							imp_cookie,imp_alexa,imp_flash,
							imp_time,imp_date,sla_id,slot_id
							) VALUES (
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,?
							)`, time.Now().Format("20060102"))

	_, err = NewManager().GetWDbMap().Exec(q,
		ca.Campaign().ID(), rh, size,
		wID, 0, appID,
		adid, copID, ca.CampaignAdID(),
		ip.String(), refer, parent,
		ca.Target(), bid, susp,
		0, alx, 0,
		ts.Unix(), ts.Format("20060102"), said, sID)
	return err
}
