package entities

import (
	"database/sql"
	"time"

	"github.com/clickyab/services/mysql"
)

//TODO: add model codegen
//TODO: add CreativesLocationsReport interface and implement it

// CreativesLocationsReport creatives_locations_reports model in database
// @Model {
//		table = creatives_locations_reports
//		primary = true, id
//		find_by = id
//		list = yes
// }
type CreativesLocationsReport struct {
	ID              int64          `json:"id" db:"id"`
	PublisherID     int64          `json:"publisher_id" db:"publisher_id"`
	PublisherDomain string         `json:"publisher_domain" db:"publisher_domain"`
	SeatID          int64          `json:"seat_id" db:"seat_id"`
	PublisherPageID int64          `json:"publisher_page_id" db:"publisher_page_id"`
	URLKey          string         `json:"url_key" db:"url_key"`
	CreativeID      int64          `json:"creative_id" db:"creative_id"`
	ActiveDays      int64          `json:"active_days" db:"active_days"`
	TotalImp        int64          `json:"total_imp" db:"total_imp"`
	TotalClicks     int64          `json:"total_clicks" db:"total_clicks"`
	TotalCTR        int64          `json:"total_ctr" db:"total_ctr"`
	YesterdayImp    int64          `json:"yesterday_imp" db:"yesterday_imp"`
	YesterdayClicks int64          `json:"yesterday_clicks" db:"yesterday_clicks"`
	YesterdayCTR    int64          `json:"yesterday_ctr" db:"yesterday_ctr"`
	TodayImp        int64          `json:"today_imp" db:"today_imp"`
	TodayClicks     int64          `json:"today_clicks" db:"today_clicks"`
	TodayCTR        int64          `json:"today_ctr" db:"today_ctr"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       mysql.NullTime `json:"updated_at" db:"updated_at"`
}

// AddAndGetCreativesLocationsReport return creative location report if exist and insert if not
func AddAndGetCreativesLocationsReport(creativeReport CreativesLocationsReport) (*CreativesLocationsReport, error) {
	fCreativeRepQ := `SELECT 
			id,
			publisher_id,
			publisher_domain,
			publisher_page_id,
			url_key,
			creative_id,
			seat_id
		FROM creatives_locations_reports
		WHERE
			publisher_id=?
			AND publisher_domain=?
			AND publisher_page_id=?
			AND url_key=?
			AND creative_id=?
			AND seat_id=?
	`
	//Important: use GetWDbMap because read db may take time to synce and fire err and finally miss impression
	err := NewManager().GetWDbMap().SelectOne(
		&creativeReport,
		fCreativeRepQ,
		creativeReport.PublisherID,
		creativeReport.PublisherDomain,
		creativeReport.PublisherPageID,
		creativeReport.URLKey,
		creativeReport.CreativeID,
		creativeReport.SeatID,
	)
	if err == nil {
		return &creativeReport, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		creativeReport.CreatedAt = time.Now()
		err = NewManager().GetWDbMap().Insert(&creativeReport)
		if err != nil {
			return nil, err
		}
	}

	return &creativeReport, nil
}
