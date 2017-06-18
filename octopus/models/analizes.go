package models

import (
	"time"

	"github.com/clickyab/services/assert"
)

const (
	// ExchangeReportTableName table name exchange
	ExchangeReportTableName = "exchange_report"
	// DemandTableName table name demand
	DemandTableName = "supplier_source_demand"
	// SupplierTableName table supplier
	SupplierTableName = "supplier_source"
	// DemandReportTableName demand report table
	DemandReportTableName = "demand_report"
	// SupplierReportTableName table supplier report
	SupplierReportTableName = "supplier_report"
)

// SupplierSourceDemand supplier_source_demand
type SupplierSourceDemand struct {
	ID              int64  `json:"id" db:"id"`
	Demand          string `json:"demand" db:"demand"`
	Supplier        string `json:"supplier" db:"supplier"`
	Source          string `json:"source" db:"source"`
	TimeID          int64  `json:"time_id" db:"time_id"`
	RequestOutCount int64  `json:"request_out_count" db:"request_out_count"`
	AdInCount       int64  `json:"ad_in_count" db:"ad_in_count"`
	ImpOutCount     int64  `json:"imp_out_count" db:"imp_out_count"`
	AdOutCount      int64  `json:"ad_out_count" db:"ad_out_count"`
	AdOutBid        int64  `json:"ad_out_bid" db:"ad_out_bid"`
	DeliverCount    int64  `json:"deliver_count" db:"deliver_count"`
	DeliverBid      int64  `json:"deliver_bid" db:"deliver_bid"`
}

// SupplierSource supplier_source
type SupplierSource struct {
	ID             int64  `json:"id" db:"id"`
	Supplier       string `json:"supplier" db:"supplier"`
	Source         string `json:"source" db:"source"`
	TimeID         int64  `json:"time_id" db:"time_id"`
	RequestInCount int64  `json:"request_in_count" db:"request_in_count"`
	ImpInCount     int64  `json:"imp_in_count" db:"imp_in_count"`
	AdOutCount     int64  `json:"ad_out_count" db:"ad_out_count"`
	DeliverCount   int64  `json:"deliver_count" db:"deliver_count"`
	DeliverBid     int64  `json:"deliver_bid" db:"deliver_bid"`
	Profit         int64  `json:"profit" db:"profit"`
}

// TimeTable TimeTable
type TimeTable struct {
	ID     int64 `json:"id" db:"id"`
	Year   int64 `json:"year" db:"year"`
	Month  int64 `json:"month" db:"month"`
	Day    int64 `json:"day" db:"day"`
	Hour   int64 `json:"hour" db:"hour"`
	JYear  int64 `json:"j_year" db:"j_year"`
	JMonth int64 `json:"j_month" db:"j_month"`
	JDay   int64 `json:"j_day" db:"j_day"`
}

// ExchangeReport exchange_report
type ExchangeReport struct {
	ID                    int64     `json:"id" db:"id"`
	Date                  time.Time `json:"target_date" db:"target_date"`
	SupplierImpressionIN  int64     `json:"supplier_impression_in" db:"supplier_impression_in"`
	SupplierImpressionOUT int64     `json:"supplier_impression_out" db:"supplier_impression_out"`
	DemandImpressionIN    int64     `json:"demand_impression_in" db:"demand_impression_in"`
	DemandImpressionOUT   int64     `json:"demand_impression_out" db:"demand_impression_out"`
	Earn                  int64     `json:"earn" db:"earn"`
	Spent                 int64     `json:"spent" db:"spent"`
	Income                int64     `json:"income" db:"income"`
}

// Parts is a multi query trick
type Parts struct {
	Query  string
	Params []interface{}
}

// DemandReport demand_report
type DemandReport struct {
	ID              int64  `json:"id" db:"id"`
	Demand          string `json:"demand" db:"demand"`
	TargetDate      int64  `json:"target_date" db:"target_date"`
	RequestOutCount int64  `json:"request_out_count" db:"request_out_count"`
	AdInCount       int64  `json:"ad_in_count" db:"ad_in_count"`
	ImpOutCount     int64  `json:"imp_out_count" db:"imp_out_count"`
	AdOutCount      int64  `json:"ad_out_count" db:"ad_out_count"`
	AdOutBid        int64  `json:"ad_out_bid" db:"ad_out_bid"`
	DeliverCount    int64  `json:"deliver_count" db:"deliver_count"`
	DeliverBid      int64  `json:"deliver_bid" db:"deliver_bid"`
	SuccessRate     int64  `json:"success_rate" db:"-"`
	DeliverRate     int64  `json:"deliver_rate" db:"-"`
	WinRate         int64  `json:"win_rate" db:"-"`
}

// SupplierReporter table
type SupplierReporter struct {
	ID             int64     `json:"id" db:"id"`
	Supplier       string    `json:"supplier" db:"supplier"`
	Date           time.Time `json:"target_date" db:"target_date"`
	ImpressionIn   int64     `json:"impression_in" db:"impression_in"`
	AdOutCount     int64     `json:"ad_out_count" db:"ad_out_count"`
	DeliveredCount int64     `json:"delivered_count" db:"delivered_count"`
	Earn           int64     `json:"earn" db:"earn"`
}

// MultiQuery is a hack to run multiple query in one transaction
func (m *Manager) MultiQuery(parts ...Parts) (err error) {
	err = m.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			assert.Nil(m.Rollback())
		} else {
			err = m.Commit()
		}
	}()

	for _, q := range parts {
		_, err = m.GetProperDBMap().Exec(q.Query, q.Params...)
		if err != nil {
			return
		}
	}

	return nil
}
