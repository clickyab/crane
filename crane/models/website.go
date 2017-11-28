package models

import "database/sql"

// website type for website
type Website struct {
	WID                int64          `json:"w_id" db:"w_id"`
	UserID             int64          `json:"u_id" db:"u_id"`
	WPubID             int64          `json:"w_pub_id" db:"w_pub_id"`
	WDomain            sql.NullString `json:"w_domain" db:"w_domain"`
	WSupplier          string         `json:"w_supplier" db:"w_supplier"`
	WName              sql.NullString `json:"w_name" db:"w_name"`
	WCategories        SharpArray     `json:"w_categories" db:"w_categories"`
	WMinBid            int64          `json:"w_minbid" db:"w_minbid"`
	WFloorCpm          sql.NullInt64  `json:"w_floor_cpm" db:"w_floor_cpm"`
	WProfileType       sql.NullInt64  `json:"w_profile_type" db:"w_profile_type"`
	WStatus            int            `json:"w_status" db:"w_status"`
	WReview            int            `json:"w_review" db:"w_review"`
	WAlexaRank         int64          `json:"w_alexarank" db:"w_alexarank"`
	WDiv               float64        `json:"w_div" db:"w_div"`
	WMobad             int            `json:"w_mobad" db:"w_mobad"`
	WNativeAd          int            `json:"w_nativead" db:"w_nativead"`
	WFatFinger         int            `json:"w_fatfinger" db:"w_fatfinger"`
	WPublishStart      int            `json:"w_publish_start" db:"w_publish_start"`
	WPublishEnd        int            `json:"w_publish_end" db:"w_publish_end"`
	WPublishCost       int            `json:"w_publish_cost" db:"w_publish_cost"`
	WPrePayment        int            `json:"w_prepayment" db:"w_prepayment"`
	WTodayCtr          float64        `json:"w_today_ctr" db:"w_today_ctr"`
	WTodayImps         int64          `json:"w_today_imps" db:"w_today_imps"`
	WTodayClicks       int64          `json:"w_today_clicks" db:"w_today_clicks"`
	WDate              int            `json:"w_date" db:"w_date"`
	WNotApprovedReason SharpArray     `json:"w_notapprovedreason" db:"w_notapprovedreason"`
	CreatedAt          sql.NullString `json:"created_at" db:"created_at"`
	UpdatedAt          sql.NullString `json:"updated_at" db:"updated_at"`
}
