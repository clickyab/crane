package entities

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"io"
	"strings"
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
	ID          int64           `json:"id" db:"id"`
	PubID       int64           `json:"publisher_id" db:"publisher_id"`
	PubDomain   string          `json:"publisher_domain" db:"publisher_domain"`
	SID         int64           `json:"seat_id" db:"seat_id"`
	PubPageID   int64           `json:"publisher_page_id" db:"publisher_page_id"`
	URLKey      string          `json:"url_key" db:"url_key"`
	CrID        int64           `json:"creative_id" db:"creative_id"`
	CrSize      int64           `json:"creative_size" db:"creative_size"`
	CrType      mysql.NullInt64 `json:"creative_type" db:"creative_type"`
	AcDays      int64           `json:"active_days" db:"active_days"`
	TotImp      int64           `json:"total_imp" db:"total_imp"`
	TotClicks   int64           `json:"total_clicks" db:"total_clicks"`
	TotCTR      int64           `json:"total_ctr" db:"total_ctr"`
	YestrImp    int64           `json:"yesterday_imp" db:"yesterday_imp"`
	YestrClicks int64           `json:"yesterday_clicks" db:"yesterday_clicks"`
	YestrCTR    int64           `json:"yesterday_ctr" db:"yesterday_ctr"`
	TodImp      int64           `json:"today_imp" db:"today_imp"`
	TodClicks   int64           `json:"today_clicks" db:"today_clicks"`
	TodCTR      int64           `json:"today_ctr" db:"today_ctr"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   mysql.NullTime  `json:"updated_at" db:"updated_at"`
}

//CreativeLocationID return location per creative ID
func (cl *CreativesLocationsReport) CreativeLocationID() int64 {
	return cl.ID
}

//SeatID return seat id of location
func (cl *CreativesLocationsReport) SeatID() int64 {
	return cl.SID
}

//PublisherPageID return publisher page id
func (cl *CreativesLocationsReport) PublisherPageID() int64 {
	return cl.PubPageID
}

//CreativeID return creative id
func (cl *CreativesLocationsReport) CreativeID() int64 {
	return cl.CrID
}

//CreativeSize return creative size of location
func (cl *CreativesLocationsReport) CreativeSize() int64 {
	return cl.CrSize
}

//ActiveDays return number of active days of location
func (cl *CreativesLocationsReport) ActiveDays() int64 {
	return cl.AcDays
}

//TotalImp return total impression od location
func (cl *CreativesLocationsReport) TotalImp() int64 {
	return cl.TotImp
}

//TotalClicks return total clicks od location
func (cl *CreativesLocationsReport) TotalClicks() int64 {
	return cl.TotClicks
}

//TotalCTR return total CTR od location
func (cl *CreativesLocationsReport) TotalCTR() int64 {
	return cl.TotCTR
}

//YesterdayImp return yesterday impression od location
func (cl *CreativesLocationsReport) YesterdayImp() int64 {
	return cl.YestrImp
}

//YesterdayClicks return yesterday clicks od location
func (cl *CreativesLocationsReport) YesterdayClicks() int64 {
	return cl.YestrClicks
}

//YesterdayCTR return yesterday CTR od location
func (cl *CreativesLocationsReport) YesterdayCTR() int64 {
	return cl.YestrCTR
}

//TodayImp return today impression od location
func (cl *CreativesLocationsReport) TodayImp() int64 {
	return cl.TodImp
}

//TodayClicks return today clicks od location
func (cl *CreativesLocationsReport) TodayClicks() int64 {
	return cl.TodClicks
}

//TodayCTR return today CTR od location
func (cl *CreativesLocationsReport) TodayCTR() int64 {
	return cl.TodCTR
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

	err := NewManager().GetRDbMap().SelectOne(
		&creativeReport,
		fCreativeRepQ,
		creativeReport.PubID,
		creativeReport.PubDomain,
		creativeReport.PubPageID,
		creativeReport.URLKey,
		creativeReport.CrID,
		creativeReport.SID,
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

// CRlocationctrLoader load all creatives statistics per locations
func CRlocationctrLoader(ids []int64) (map[int64]map[string]kv.Serializable, error) {
	crlocationctr := make(map[int64]map[string]kv.Serializable)
	// return crlocationctr, nil // Uncomment this line after first time in DEV mode
	if len(ids) == 0 {
		return crlocationctr, fmt.Errorf("no creatives id")
	}

	// yesterday, _ := strconv.Atoi(time.Now().AddDate(0, 0, -1).Format("20060102"))
	p := make([]string, len(ids))
	params := make([]interface{}, len(ids))
	c := 0
	for i := range ids {
		p[c] = "?"
		params[c] = ids[i]
		c++
	}
	cond := fmt.Sprintf("creative_id IN(%s)", strings.Join(p, ","))

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
			WHERE %s
			LIMIT %d, %d`,
			cond,
			j,
			j+cnt,
		)

		var res []CreativesLocationsReport
		if _, err := NewManager().GetRDbMap().Select(&res, q, params...); err != nil {
			logrus.Warn(err)
			return crlocationctr, err
		}

		if len(res) == 0 {
			break
		}

		for i := range res {
			key := GenCRPerLocationPoolKey(
				res[i].PubDomain,
				res[i].PubPageID,
				res[i].SID,
				res[i].CrID,
				res[i].CrSize,
			)
			crlocationctr[res[i].CrID][key] = &res[i]
		}
	}

	logrus.Debugf("Load %d creatives statistics per locations", len(crlocationctr))

	return crlocationctr, nil
}

// Encode is the encode function for serialize object in io writer
func (cl CreativesLocationsReport) Encode(w io.Writer) error {
	g := gob.NewEncoder(w)
	return g.Encode(cl)
}

// Decode try to decode object from io reader
func (cl CreativesLocationsReport) Decode(r io.Reader) error {
	g := gob.NewDecoder(r)
	return g.Decode(cl)
}

// GenCRPerLocationPoolKey generate cache key for pool
func GenCRPerLocationPoolKey(publisherDomain string, pageID, seatID, creativeID, creativeSize int64) string {
	return fmt.Sprintf(
		"pubdom%s_page%d_seat%d_creative%d_size%d",
		publisherDomain,
		pageID,
		seatID,
		creativeID,
		creativeSize,
	)
}
