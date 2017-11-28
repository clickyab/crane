package models

import (
	"database/sql"
	"time"
)

// app is the applications structure
type App struct {
	ID                   int64          `db:"app_id"`
	UserID               int64          `db:"u_id"`
	AppToken             string         `db:"app_token"`
	AppName              string         `db:"app_name"`
	EnAppName            string         `db:"en_app_name"`
	AppPackage           string         `db:"app_package"`
	AppSupplier          string         `db:"app_supplier"`
	AmID                 int            `db:"am_id"`
	MinBID               int64          `db:"app_minbid"`
	AppFloorCPM          sql.NullInt64  `db:"app_floor_cpm"`
	AppDIV               float64        `db:"app_div"`
	AppStatus            int            `db:"app_status"`
	AppReview            int            `db:"app_review"`
	AppTodayCTR          int64          `db:"app_today_ctr"`
	AppTodayIMPs         int64          `db:"app_today_imps"`
	AppTodayClicks       int64          `db:"app_today_clicks"`
	AppDate              int            `db:"app_date"`
	Appcat               SharpArray     `db:"app_cat"`
	AppNotApprovedReason sql.NullString `db:"app_notapprovedreason"`
	AppFatFinger         sql.NullBool   `db:"app_fatfinger"`
	CreatedAt            time.Time      `db:"created_at"`
	UpdatedAt            time.Time      `db:"updated_at"`

	AppPrepayment  int `db:"app_prepayment"`
	AppPublishCost int `db:"app_publish_cost"`
}

// phoneData is the phone data united in one structure for filtering
type PhoneData struct {
	Brand   string
	BrandID int64
	// Model     string
	// ModelID   int64
	Carrier   string
	CarrierID int64
	// Lang      string
	// LangID    int64
	Network   string
	NetworkID int64
}
