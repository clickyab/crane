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
func AddImpression(rh, ref, par, spid, copID string, size, susp int, adid int64, pub entity.Publisher, ip net.IP,
	bid float64, alexa bool, ts time.Time, typ entity.RequestType, cpm, scpm float64) error {
	var err error
	impCPM := cpm
	if scpm != 0 {
		impCPM = scpm
	}
	sDiffCPM := sql.NullFloat64{Valid: scpm != 0, Float64: cpm - scpm}
	wID := sql.NullInt64{}
	appID := sql.NullInt64{}
	refer := sql.NullString{}
	parent := sql.NullString{}
	if typ == entity.RequestTypeDemand {
		// TODO : check for publisher type in demand too
		wID = sql.NullInt64{Valid: err == nil, Int64: pub.ID()}
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
							imp_time,imp_date,sla_id,
							slot_id, s_name, s_diff_cpm,
							imp_cpm
							) VALUES (
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?
							)`, time.Now().Format("20060102"))

	_, err = NewManager().GetWDbMap().Exec(q,
		ca.Campaign().ID(), rh, size,
		wID, 0, appID,
		adid, copID, ca.CampaignAdID(),
		ip.String(), refer, parent,
		ca.Target(), bid, susp,
		0, alx, 0,
		ts.Unix(), ts.Format("20060102"), said,
		sID, pub.Supplier().Name(), sDiffCPM,
		impCPM)
	return err
}
