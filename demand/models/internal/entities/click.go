package entities

import (
	"net"
	"time"

	"database/sql"

	"clickyab.com/crane/demand/entity"
)

// Click fill structure for Click
type Click struct {
	reservedHash string
	winnerBid    float64
	adID         int64
	webSiteID    int64
	appID        int64
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
func FillClickData(supplier, rh, slotPubID, ref, parent, os, copID string, susp, size int, fast, adID int64, bid float64, ip net.IP, ts time.Time, pubID int64) (*Click, error) {
	// find slot
	slotID, err := FindWebSlotID(slotPubID, pubID, size)
	if err != nil {
		return nil, err
	}

	//find ad
	ad, err := GetAd(adID)
	if err != nil {
		return nil, err
	}
	//find slot ad
	slotAdID, err := FindSlotAd(slotID, ad.ID())
	if err != nil {
		return nil, err
	}
	return &Click{
		reservedHash: rh,
		winnerBid:    bid,
		webSiteID:    pubID,
		appID:        0, //TODO should be filled after App implemented
		campaignID:   ad.Campaign().ID(),
		campaignAdID: ad.CampaignAdID(),
		slotID:       slotID,
		slotAdID:     slotAdID,
		copID:        copID,
		impID:        0, //TODO not sure about that
		status:       susp,
		ip:           ip.String(),
		referrer:     ref,
		parent:       parent,
		os:           os,
		adSize:       size,
		time:         ts.Unix(),
		date:         ts.Format("20060102"),
		fast:         fast, //TODO fix after rebase,
		typ:          entity.RequestTypeDemand,
		adID:         ad.ID(),
		supplier:     supplier,
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

	wid := sql.NullInt64{Valid: c.webSiteID != 0, Int64: c.webSiteID}
	referrer := sql.NullString{Valid: c.referrer != "", String: c.referrer}
	parent := sql.NullString{Valid: c.parent != "", String: c.parent}

	_, err := NewManager().GetWDbMap().Exec(q,
		c.reservedHash, c.winnerBid, wid, 0, 0, c.campaignID,
		c.campaignAdID, c.slotID, c.slotAdID, c.adID, c.copID, 0, c.status,
		c.ip, referrer, parent, c.fast, c.os, c.time, c.date, c.adSize, c.supplier)

	return err
}
