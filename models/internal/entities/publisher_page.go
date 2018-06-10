package entities

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/kv"
	"github.com/sirupsen/logrus"

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
			AND publisher_domain=?
	`
	var publisherPage PublisherPage

	err = NewManager().GetRDbMap().SelectOne(
		&publisherPage,
		fPubPageQ,
		urlKey,
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
	url, err := NormalizeURL(url)
	assert.Nil(err)

	hasher := md5.New()
	_, err = hasher.Write([]byte(url))
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

// PagesLoader load all publisher pages
func PagesLoader() func(ctx context.Context) (map[string]kv.Serializable, error) {
	return func(ctx context.Context) (map[string]kv.Serializable, error) {
		pages := make(map[string]kv.Serializable)
		//TODO: should comment after fix it
		return pages, nil // Uncomment this line after first time in DEV mode

		yesterday, _ := strconv.Atoi(time.Now().AddDate(0, 0, -1).Format("20060102"))
		const cnt = 10000
		for j := 0; ; j = j + cnt {
			q := fmt.Sprintf(`SELECT 
					id,
					publisher_id,
					publisher_domain,
					kind,
					url,
					url_key,
					active_days,
					avg_daily_imp,
					avg_daily_clicks,
					today_imp,
					today_clicks,
					today_ctr
				FROM publisher_pages
				WHERE updated_at IS NULL OR updated_at>?
				LIMIT %d, %d`,
				j,
				j+cnt,
			)

			var res []PublisherPage
			if _, err := NewManager().GetRDbMap().Select(&res, q, yesterday); err != nil {
				logrus.Warn(err)
				return nil, err
			}

			if len(res) == 0 {
				break
			}

			for i := range res {
				key := GenPagePoolKey(
					res[i].PublisherDomain,
					res[i].URLKey,
				)
				pages[key] = &res[i]
			}
		}

		logrus.Debugf("Load %d publisher pages", len(pages))

		return pages, nil
	}
}

// GenPagePoolKey generate cache key for pool
func GenPagePoolKey(publisherDomain, URLKey string) string {
	return fmt.Sprintf(
		"pubdo%s_url%",
		publisherDomain,
		URLKey,
	)
}

// Encode is the encode function for serialize object in io writer
func (p PublisherPage) Encode(w io.Writer) error {
	g := gob.NewEncoder(w)
	return g.Encode(p)
}

// Decode try to decode object from io reader
func (p PublisherPage) Decode(r io.Reader) error {
	g := gob.NewDecoder(r)
	return g.Decode(p)
}
