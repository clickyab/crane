package entities

import (
	"database/sql"
	"strconv"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/workers/models"
)

// Click fill structure for Click
type Click struct {
	reservedHash string
	winnerBid    float64
	adID         int32
	webSiteID    sql.NullInt64
	appID        sql.NullInt64
	campaignID   int32
	campaignAdID int32
	slotID       int64
	slotAdID     int32
	copID        string
	impID        int64
	tv           bool
	status       int
	ip           string
	referrer     string
	parent       string
	fast         int64
	os           string
	time         int64
	date         string
	supplier     string
	adSize       int32
	typ          entity.InputType
}

// FindImpFromClickByImpID return impression by impression id
func FindImpFromClickByImpID(imp int64) (*Impression, error) {
	q := `SELECT w_id,app_id,wp_id,ca_id,ad_id,cop_id,cp_id,slot_id,imp_id, reserved_hash
				FROM  clicks WHERE imp_id = ?`
	var x = &Impression{}
	err := NewManager().GetRDbMap().SelectOne(x, q, imp)

	return x, err
}

// FindImpFromClickByRH return impression by reserved hash
func FindImpFromClickByRH(rh string) (*Impression, error) {
	q := `SELECT w_id,app_id,wp_id,ca_id,ad_id,cop_id,cp_id,slot_id,imp_id, reserved_hash
				FROM  clicks WHERE reserved_hash = ?`
	var x = &Impression{}
	err := NewManager().GetRDbMap().SelectOne(x, q, rh)

	return x, err
}

// FillClickData try to fill Click structure
func FillClickData(p entity.Publisher, m models.Impression, s models.Seat, os entity.OS, fast int64, tv bool) (*Click, error) {

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
		tv:           tv,
		adSize:       s.AdSize,
		time:         m.Timestamp.Unix(),
		date:         m.Timestamp.Format("20060102"),
		fast:         fast,
		typ:          entity.InputTypeDemand,
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
	copString := c.copID
	if len(c.copID) > 10 {
		copString = copString[:10]
	}
	copID, _ := strconv.ParseInt(copString, 16, 64)

	referrer := sql.NullString{Valid: c.referrer != "", String: c.referrer}
	parent := sql.NullString{Valid: c.parent != "", String: c.parent}

	res, err := NewManager().GetWDbMap().Exec(q,
		c.reservedHash, c.winnerBid, c.webSiteID, c.appID, 0, c.campaignID,
		c.campaignAdID, c.slotID, c.slotAdID, c.adID, copID, 0, c.status,
		c.ip, referrer, parent, c.fast, c.os, c.time, c.date, c.adSize, c.supplier)

	if err != nil {
		return err
	}

	// insert true view here
	if c.tv {
		clickID, err := res.LastInsertId()
		if err != nil {
			return err
		}
		err = insertTrueView(clickID)
		if err != nil {
			return err
		}

	}

	return err
}

func insertTrueView(cID int64) error {
	q := `INSERT INTO trueview (tv_click_id) VALUES (?)`
	_, err := NewManager().GetWDbMap().Exec(q, cID)
	return err
}
