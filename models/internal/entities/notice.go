package entities

import (
	"fmt"

	"time"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/workers/models"
)

// AddNotice get impression from job abd insert it into notice table
func AddNotice(p entity.Publisher, m models.Impression, s models.Seat) error {
	ca, err := GetAd(s.AdID)
	if err != nil {
		return err
	}
	q := fmt.Sprintf("INSERT INTO win_requests (hash, supplier, publisher_id, campaign_id, creative_id, cpc, cpm, `type`, created_at) " +
		"VALUES ( ?,?,?,?,?,?,?,?,? )")
	_, err = NewManager().GetWDbMap().Exec(q, s.ReserveHash, m.Supplier, p.ID(), ca.ID(), s.AdID, s.WinnerBID, s.CPM, s.Type.String(), time.Now())
	return err
}
