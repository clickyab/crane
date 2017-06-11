package models

import (
	"time"

	"fmt"

	"clickyab.com/exchange/services/assert"
)

// fetchDemand select demand side
func fetchDemand(start int64, end int64) *Exchange {
	ex := Exchange{}
	q := fmt.Sprintf(`SELECT SUM(imp_in_count) AS demand_impression_in,
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
func fetchSupplier(start int64, end int64) *Exchange {
	ex := Exchange{}
	q := fmt.Sprintf(`SELECT SUM(request_in_count) AS supplier_impression_in,
	SUM(deliver_count) AS supplier_impression_out,
	SUM(deliver_bid) AS spent
	FROM %s
	WHERE time_id >= ?
	AND time_id <= ?`, SuplierTableName)
	m := NewManager()
	_, err := m.GetRDbMap().Select(ex, q, start, end)
	assert.Nil(err)
	return &ex
}

// ExchangeReport cron worker report exchange
func ExchangeReport(date time.Time) {
	start, end := factTableYesterdayID(date)
	dem := fetchDemand(start, end)
	sup := fetchSupplier(start, end)
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
				`, ExchangeTableName)
	m := NewManager()
	_, err := m.GetRDbMap().Exec(q, date, sup.SupplierImpressionIN,
		sup.SupplierImpressionOUT, dem.DemandImpressionIN, dem.DemandImpressionOUT,
		sup.Earn, dem.Spent, sup.Earn-sup.Spent)
	assert.Nil(err)
}

// FactTableYesterdayID is a helper function to get the fact table for yesterday id from time
func factTableYesterdayID(tm time.Time) (int64, int64) {
	y, m, d := tm.Date()
	from := time.Date(y, m, d, 0, 0, 1, 0, time.UTC)
	to := time.Date(y, m, d, 23, 59, 59, 0, time.UTC)
	return FactTableID(from), FactTableID(to)
}
