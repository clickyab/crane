package entities

import (
	"database/sql"
	"time"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/workers/models"
	"github.com/clickyab/services/mysql"
)

//TODO: add model codegen
//TODO: fix seat interface and implement it

// Seat seats model in database
// @Model {
//		table = seats
//		primary = true, id
//		find_by = id
//		list = yes
// }
type Seat struct {
	ID              int64                `json:"id" db:"id"`
	SlotID          int64                `json:"slot_id" db:"slot_id"`
	SupplierName    string               `json:"supplier_name" db:"supplier_name"`
	PublisherID     int64                `json:"publisher_id" db:"publisher_id"`
	PublisherDomain string               `json:"publisher_domain" db:"publisher_domain"`
	Kind            entity.PublisherType `json:"kind" db:"kind"`
	ActiveDays      int64                `json:"active_days" db:"active_days"`
	AvgDailyImp     int64                `json:"avg_daily_imp" db:"avg_daily_imp"`
	AvgDailyClicks  int64                `json:"avg_daily_clicks" db:"avg_daily_clicks"`
	TodayImp        int64                `json:"today_imp" db:"today_imp"`
	TodayClicks     int64                `json:"today_clicks" db:"today_clicks"`
	TodayCTR        int64                `json:"today_ctr" db:"today_ctr"`
	CreatedAt       time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt       mysql.NullTime       `json:"updated_at" db:"updated_at"`
}

// AddAndGetSeat return seat if exist and insert if not
func AddAndGetSeat(m models.Impression, slID int64) (*Seat, error) {
	var seat Seat

	fSeatQ := `SELECT 
			id,
			slot_id,
			supplier_name,
			publisher_id,
			publisher_domain
		FROM seats
		WHERE
			slot_id=?
			AND supplier_name=?
			AND publisher_id=?
			AND publisher_domain=?
	`

	//Important: use GetWDbMap because read db may take time to synce and fire err and finally miss impression
	err := NewManager().GetWDbMap().SelectOne(
		&seat,
		fSeatQ,
		slID,
		m.Supplier,
		m.PublisherID,
		m.Publisher,
	)
	if err == nil {
		return &seat, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		seat.SlotID = slID
		seat.SupplierName = m.Supplier
		seat.PublisherID = m.PublisherID
		seat.PublisherDomain = m.Publisher
		seat.Kind = m.PublisherType
		seat.CreatedAt = time.Now()

		err = NewManager().GetWDbMap().Insert(&seat)
		if err != nil {
			return nil, err
		}
	}

	return &seat, nil
}
