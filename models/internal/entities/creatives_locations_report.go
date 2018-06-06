package entities

import (
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/mysql"
	"github.com/sirupsen/logrus"
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
	ID              int64           `json:"id" db:"id"`
	PublisherID     int64           `json:"publisher_id" db:"publisher_id"`
	PublisherDomain string          `json:"publisher_domain" db:"publisher_domain"`
	SeatID          int64           `json:"seat_id" db:"seat_id"`
	PublisherPageID int64           `json:"publisher_page_id" db:"publisher_page_id"`
	URLKey          string          `json:"url_key" db:"url_key"`
	CreativeID      int64           `json:"creative_id" db:"creative_id"`
	CreativeSize    int64           `json:"creative_size" db:"creative_size"`
	CreativeType    mysql.NullInt64 `json:"creative_type" db:"creative_type"`
	ActiveDays      int64           `json:"active_days" db:"active_days"`
	TotalImp        int64           `json:"total_imp" db:"total_imp"`
	TotalClicks     int64           `json:"total_clicks" db:"total_clicks"`
	TotalCTR        int64           `json:"total_ctr" db:"total_ctr"`
	YesterdayImp    int64           `json:"yesterday_imp" db:"yesterday_imp"`
	YesterdayClicks int64           `json:"yesterday_clicks" db:"yesterday_clicks"`
	YesterdayCTR    int64           `json:"yesterday_ctr" db:"yesterday_ctr"`
	TodayImp        int64           `json:"today_imp" db:"today_imp"`
	TodayClicks     int64           `json:"today_clicks" db:"today_clicks"`
	TodayCTR        int64           `json:"today_ctr" db:"today_ctr"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       mysql.NullTime  `json:"updated_at" db:"updated_at"`
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
			creative_size,
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

// CRPerLocationsLoader load all creatives statistics per locations
func CRPerLocationsLoader() func(ctx context.Context) (map[string]kv.Serializable, error) {
	return func(ctx context.Context) (map[string]kv.Serializable, error) {
		crPerLocations := make(map[string]kv.Serializable)
		// return b, nil // Uncomment this line after first time in DEV mode

		yesterday, _ := strconv.Atoi(time.Now().AddDate(0, 0, -1).Format("20060102"))
		const cnt = 10000
		for j := 0; ; j = j + cnt {
			q := fmt.Sprintf(`SELECT 
					id,
					publisher_id,
					publisher_domain,
					publisher_page_id,
					url_key,
					creative_id,
					creative_size,
					seat_id,
					active_days,
					total_imp,
					total_clicks,
					total_ctr,
					yesterday_imp,
					yesterday_clicks,
					yesterday_ctr,
					today_imp,
					today_clicks,
					today_ctr
				FROM creatives_locations_reports
				WHERE updated_at IS NULL OR updated_at>?
				LIMIT %d, %d`,
				j,
				j+cnt,
			)

			var res []CreativesLocationsReport
			if _, err := NewManager().GetRDbMap().Select(&res, q, yesterday); err != nil {
				logrus.Warn(err)
				return nil, err
			}

			if len(res) == 0 {
				break
			}

			for i := range res {
				key := GenCRPerLocationPoolKey(
					res[i].PublisherDomain,
					res[i].PublisherPageID,
					res[i].SeatID,
					res[i].CreativeID,
				)
				crPerLocations[key] = &res[i]
			}
		}

		logrus.Debugf("Load %d creatives statistics per locations", len(crPerLocations))

		return crPerLocations, nil
	}
}

// GenCRPerLocationPoolKey generate cache key for pool
func GenCRPerLocationPoolKey(publisherDomain string, pageID, seatID, creativeID int64) string {
	return fmt.Sprintf(
		"pubdom%s_page%d_seat%d_creative%d",
		publisherDomain,
		pageID,
		seatID,
		creativeID,
	)
}

// Encode is the encode function for serialize object in io writer
func (cr CreativesLocationsReport) Encode(w io.Writer) error {
	g := gob.NewEncoder(w)
	return g.Encode(cr)
}

// Decode try to decode object from io reader
func (cr CreativesLocationsReport) Decode(r io.Reader) error {
	g := gob.NewDecoder(r)
	return g.Decode(cr)
}
