package models

import (
	"time"

	"fmt"

	"github.com/clickyab/services/assert"
)

// fetchDemand select demand side
func fetchDemand(start int64, end int64) *ExchangeReport {
	ex := ExchangeReport{}
	q := fmt.Sprintf(`SELECT SUM(ad_out_count) AS demand_impression_in,
	SUM(imp_out_count) AS demand_impression_out,
	SUM(deliver_bid) AS earn
	FROM %s
	WHERE time_id >= ?
	AND time_id <= ?`, DemandTableName)
	m := NewManager()
	_, err := m.GetRDbMap().Select(ex, q, start, end)
	assert.Nil(err)
	return &ex
}

// fetchSupplier select demand side
func fetchSupplier(start int64, end int64) *ExchangeReport {
	ex := ExchangeReport{}
	q := fmt.Sprintf(`SELECT SUM(request_in_count) AS supplier_impression_in,
	SUM(deliver_count) AS supplier_impression_out,
	SUM(deliver_bid) AS spent
	FROM %s
	WHERE time_id >= ?
	AND time_id <= ?`, SupplierTableName)
	m := NewManager()
	_, err := m.GetRDbMap().Select(ex, q, start, end)
	assert.Nil(err)
	return &ex
}

// updateExchangeReport will update demand report (inclusive)
func updateExchangeReport(t time.Time) {
	td := t.Format("2006-01-02")
	from, to := factTableRange(t)
	dem := fetchDemand(from, to)
	sup := fetchSupplier(from, to)
	q := fmt.Sprintf(`INSERT INTO %s
				(target_date,
				supplier_impression_in,
				supplier_impression_out,
				demand_impression_in,
				demand_impression_out,
				earn,
				spent,
				income)
				VALUES(?,?,?,?,?,?,?,?)
				ON DUPLICATE KEY UPDATE
				supplier_impression_in = VALUES(supplier_impression_in),
				supplier_impression_out = VALUES(supplier_impression_out),
				demand_impression_in = VALUES(demand_impression_in),
				demand_impression_out = VALUES(demand_impression_out),
				earn = VALUES(earn),
				spent = VALUES(spent),
				income = VALUES(income)
				`, ExchangeReportTableName)
	m := NewManager()
	_, err := m.GetRDbMap().Exec(q, td, sup.SupplierImpressionIN,
		sup.SupplierImpressionOUT, dem.DemandImpressionIN, dem.DemandImpressionOUT,
		sup.Earn, dem.Spent, sup.Earn-sup.Spent)
	assert.Nil(err)
}

// UpdateExchangeReportRange cron worker report exchange
func UpdateExchangeReportRange(from time.Time, to time.Time) {
	if from.Unix() > to.Unix() {
		from, to = to, from
	}
	to = to.Add(24 * time.Hour)
	for from.Unix() < to.Unix() {
		updateExchangeReport(from)
		from = from.Add(time.Hour * 24)
	}

}
