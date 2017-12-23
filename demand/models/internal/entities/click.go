package entities

import (
	"database/sql"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/workers/models"
)

// Click fill structure for Click
type Click struct {
	reservedHash string
	winnerBid    float64
	adID         int64
	webSiteID    sql.NullInt64
	appID        sql.NullInt64
	campaignID   int64
	campaignAdID int64
	slotID       int64
	slotAdID     int64
	copID        string
	impID        int64
	status       int
	ip           string
	referrer     string
	parent       string
	fast         int64
	os           string
	time         int64
	date         string
	supplier     string
	adSize       int
	typ          entity.RequestType
}

// FillClickData try to fill Click structure
func FillClickData(p entity.Publisher, m models.Impression, s models.Seat, os entity.OS, fast int64) (*Click, error) {

	var err error

	wID := sql.NullInt64{}
	appID := sql.NullInt64{}

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
		return nil, err
	}

	//find ad
	ad, err := GetAd(s.AdID)
	if err != nil {
		return nil, err
	}

	return &Click{
		reservedHash: s.ReserveHash,
		winnerBid:    s.WinnerBID,
		webSiteID:    wID,
		appID:        appID,
		campaignID:   ad.Campaign().ID(),
		campaignAdID: ad.CampaignAdID(),
		slotID:       sID,
		slotAdID:     s.AdID,
		copID:        m.CopID,
		impID:        0, //TODO not sure about that
		status:       m.Suspicious,
		ip:           m.IP.String(),
		referrer:     m.Referrer,
		parent:       m.ParentURL,
		os:           os.Name,
		adSize:       s.AdSize,
		time:         m.Timestamp.Unix(),
		date:         m.Timestamp.Format("20060102"),
		fast:         fast, //TODO fix after rebase,
		typ:          entity.RequestTypeDemand,
		adID:         ad.ID(),
		supplier:     p.Supplier().Name(),
	}, nil

}

// InsertClick try to inset Click
func InsertClick(c *Click) error {
	q := `INSERT INTO clicks (reserved_hash,c_winnerbid,
	w_id,
	app_id,
	wp_id,
	cp_id,
	ca_id,
	slot_id,
	sla_id,
	ad_id,
	cop_id,
	imp_id,
	c_status,
	c_ip,
	c_referaddress,
	c_parenturl,
	c_fast,
	c_os,
	c_time,
	c_date,ad_size,c_supplier) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	referrer := sql.NullString{Valid: c.referrer != "", String: c.referrer}
	parent := sql.NullString{Valid: c.parent != "", String: c.parent}

	_, err := NewManager().GetWDbMap().Exec(q,
		c.reservedHash, c.winnerBid, c.webSiteID, c.appID, 0, c.campaignID,
		c.campaignAdID, c.slotID, c.slotAdID, c.adID, c.copID, 0, c.status,
		c.ip, referrer, parent, c.fast, c.os, c.time, c.date, c.adSize, c.supplier)

	return err
}
