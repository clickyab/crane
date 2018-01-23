package entities

import (
	"database/sql"
	"fmt"
	"time"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/workers/models"
)

// AddImpression insert new impression to daily table
// TODO : multiple insert per one query
func AddImpression(p entity.Publisher, m models.Impression, s models.Seat) error {
	var err error
	impCPM := s.CPM
	if s.SCPM != 0 {
		impCPM = s.SCPM
	}
	sDiffCPM := sql.NullFloat64{Valid: s.SCPM != 0, Float64: s.CPM - s.SCPM}
	wID := sql.NullInt64{}
	appID := sql.NullInt64{}
	refer := sql.NullString{Valid: m.Referrer != "", String: m.Referrer}
	parent := sql.NullString{Valid: m.ParentURL != "", String: m.ParentURL}

	if p.Type() == entity.PublisherTypeWeb {
		wID = sql.NullInt64{Valid: p.ID() != 0, Int64: p.ID()}

	} else if p.Type() == entity.PublisherTypeApp {
		appID = sql.NullInt64{Valid: p.ID() != 0, Int64: p.ID()}
	} else {
		panic("mismatch impression and publisher type")
	}

	var sID int64

	// find slot id
	if wID.Valid {
		sID, err = FindWebSlotID(s.SlotPublicID, wID.Int64, s.AdSize)
	} else if appID.Valid {
		sID, err = FindAppSlotID(s.SlotPublicID, appID.Int64, s.AdSize)
	}
	if err != nil {
		return err
	}

	// find slot ad
	said, err := FindSlotAd(sID, s.AdID)
	if err != nil {
		return err
	}
	ca, err := GetAd(s.AdID)
	if err != nil {
		return err
	}
	var alx int
	if m.Alexa {
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
		ca.Campaign().ID(), s.ReserveHash, s.AdSize,
		wID, 0, appID,
		s.AdID, m.CopID, ca.CampaignAdID(),
		m.IP.String(), refer, parent,
		ca.TargetURL(), s.WinnerBID, m.Suspicious,
		0, alx, 0,
		m.Timestamp.Unix(), m.Timestamp.Format("20060102"), said,
		sID, p.Supplier().Name(), sDiffCPM,
		impCPM)
	return err
}
