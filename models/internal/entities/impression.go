package entities

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/workers/models"
	"github.com/clickyab/services/mysql"
)

// Impression model for database
type Impression struct {
	WID          mysql.NullInt64  `db:"w_id"`
	AppID        mysql.NullInt64  `db:"app_id"`
	WpID         mysql.NullInt64  `db:"wp_id"`
	CaID         mysql.NullInt64  `db:"ca_id"`
	AdID         mysql.NullInt64  `db:"ad_id"`
	CopID        mysql.NullInt64  `db:"cop_id"`
	CpID         mysql.NullInt64  `db:"cp_id"`
	SlotID       mysql.NullInt64  `db:"slot_id"`
	ImpID        mysql.NullInt64  `db:"imp_id"`
	ReservedHash mysql.NullString `db:"reserved_hash"`
}

func impTableName(t time.Time) string {
	return fmt.Sprintf("impressions%s", t.Format("20060102"))
}

// FindImpressionByID return impression by impression id
func FindImpressionByID(impid int64, t time.Time) (*Impression, error) {

	q := fmt.Sprintf(`SELECT w_id,app_id,wp_id,ca_id,ad_id,cop_id,cp_id,slot_id,imp_id, reserved_hash
				FROM  %s WHERE imp_id = ?`, impTableName(t))
	var x = &Impression{}
	err := NewManager().GetWDbMap().SelectOne(x, q, impid)
	return x, err
}

// FindImpressionByRH return impression by reserved hash
func FindImpressionByRH(rh string, t time.Time) (*Impression, error) {
	q := fmt.Sprintf(`SELECT w_id,app_id,wp_id,ca_id,ad_id,cop_id,cp_id,slot_id,imp_id, reserved_hash
				FROM  %s WHERE reserved_hash = ?`, impTableName(t))
	var x = &Impression{}
	err := NewManager().GetWDbMap().SelectOne(x, q, rh)

	return x, err
}

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
	copString := m.CopID
	if len(m.CopID) > 10 {
		copString = copString[:10]
	}
	copID, _ := strconv.ParseInt(copString, 16, 64)
	q := fmt.Sprintf(`INSERT INTO impressions%s (
							cp_id,reserved_hash,ad_size,
							w_id,wp_id,app_id,
							ad_id,cop_id,ca_id,
							imp_ipaddress,imp_referaddress,imp_parenturl,
							imp_url,imp_winnerbid,imp_status,
							imp_cookie,imp_flash,
							imp_time,imp_date,sla_id,
							slot_id, s_name, s_diff_cpm,
							imp_cpm,imp_final_cpm
							) VALUES (
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?
							)`, time.Now().Format("20060102"))

	_, err = NewManager().GetWDbMap().Exec(q,
		ca.Campaign().ID(), s.ReserveHash, s.AdSize,
		wID, 0, appID,
		s.AdID, copID, ca.CampaignAdID(),
		m.IP.String(), refer, parent,
		ca.TargetURL(), s.WinnerBID, m.Suspicious,
		0, 0,
		m.Timestamp.Unix(), m.Timestamp.Format("20060102"), said,
		sID, p.Supplier().Name(), sDiffCPM,
		impCPM, impCPM*float64(p.Supplier().Share())/100)
	if err != nil {
		return err
	}

	seat, err := AddAndGetSeat(m, sID)
	if err != nil {
		return err
	}

	pubPage, err := AddAndGetPublisherPage(m)
	if err != nil {
		return err
	}

	crReport := CreativesLocationsReport{
		PublisherID:     m.PublisherID,
		PublisherDomain: m.Publisher,
		PublisherPageID: pubPage.ID,
		URLKey:          pubPage.URLKey,
		CreativeID:      s.AdID,
		SeatID:          seat.ID,
	}
	_, err = AddAndGetCreativesLocationsReport(crReport)

	return err
}
