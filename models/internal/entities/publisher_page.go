package entities

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/clickyab/services/assert"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/workers/models"
	"github.com/PuerkitoBio/purell"
	"github.com/clickyab/services/mysql"
)

//TODO: add model codegen
//TODO: add PublisherPage interface and implement it

// PublisherPage publisher_pages model in database
// @Model {
//		table = publisher_pages
//		primary = true, id
//		find_by = id
//		list = yes
// }
type PublisherPage struct {
	ID              int64                `json:"id" db:"id"`
	PublisherID     int64                `json:"publisher_id" db:"publisher_id"`
	PublisherDomain string               `json:"publisher_domain" db:"publisher_domain"`
	Kind            entity.PublisherType `json:"kind" db:"kind"`
	URL             string               `json:"url" db:"url"`
	URLKey          string               `json:"url_key" db:"url_key"`
	ActiveDays      int64                `json:"active_days" db:"active_days"`
	AvgDailyImp     int64                `json:"avg_daily_imp" db:"avg_daily_imp"`
	AvgDailyClicks  int64                `json:"avg_daily_clicks" db:"avg_daily_clicks"`
	TodayImp        int64                `json:"today_imp" db:"today_imp"`
	TodayClicks     int64                `json:"today_clicks" db:"today_clicks"`
	TodayCTR        int64                `json:"today_ctr" db:"today_ctr"`
	CreatedAt       time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt       mysql.NullTime       `json:"updated_at" db:"updated_at"`
}

// AddAndGetPublisherPage return publisher_pages if exist and insert if not
func AddAndGetPublisherPage(m models.Impression) (*PublisherPage, error) {
	url, err := NormalizeURL(m.ParentURL)
	if err != nil {
		return nil, err
	}
	urlKey := GenerateURLKey(url)

	fPubPageQ := `SELECT 
			id,
			publisher_id,
			publisher_domain,
			kind,
			url,
			url_key
		FROM publisher_pages
		WHERE
			url_key=?
			AND publisher_id=?
			AND publisher_domain=?
	`
	var publisherPage PublisherPage
	//Important: use GetWDbMap because read db may take time to synce and fire err and finally miss impression
	err = NewManager().GetWDbMap().SelectOne(
		&publisherPage,
		fPubPageQ,
		urlKey,
		m.PublisherID,
		m.Publisher,
	)
	if err == nil {
		return &publisherPage, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		publisherPage.PublisherID = m.PublisherID
		publisherPage.PublisherDomain = m.Publisher
		publisherPage.Kind = m.PublisherType
		publisherPage.URL = url
		publisherPage.URLKey = urlKey
		publisherPage.CreatedAt = time.Now()

		err = NewManager().GetWDbMap().Insert(&publisherPage)
		if err != nil {
			return nil, err
		}
	}

	return &publisherPage, nil
}

// GenerateURLKey generate url key
func GenerateURLKey(url string) string {
	hasher := md5.New()
	_, err := hasher.Write([]byte(url))
	assert.Nil(err)
	return hex.EncodeToString(hasher.Sum(nil))
}

//NormalizeURL normalize url
func NormalizeURL(url string) (string, error) {
	return purell.NormalizeURLString(
		url,
		purell.FlagLowercaseScheme|
			purell.FlagLowercaseHost|
			purell.FlagUppercaseEscapes|
			purell.FlagRemoveDefaultPort|
			purell.FlagRemoveEmptyQuerySeparator|
			purell.FlagRemoveTrailingSlash|
			purell.FlagRemoveFragment|
			purell.FlagRemoveDuplicateSlashes|
			purell.FlagSortQuery|
			purell.FlagForceHTTP,
	)
}
