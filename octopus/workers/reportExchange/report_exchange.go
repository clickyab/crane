package reportExchange

import (
	"time"

	"fmt"

	"clickyab.com/exchange/octopus/models"
	"clickyab.com/exchange/octopus/workers/internal/datamodels"
	"clickyab.com/exchange/services/assert"
)

// FetchDemand select demand side
func FetchDemand(start int64, end int64) *models.Exchange {
	ex := models.Exchange{}
	q := fmt.Sprintf(`SELECT SUM(imp_in_count) AS demand_impression_in,
	SUM(imp_out_count) AS demand_impression_out,
	SUM(deliver_bid) AS earn
	FROM %s
	WHERE time_id >= ?
	AND time_id <= ?`, models.DemandTableName)
	m := models.NewManager()
	_, err := m.GetRDbMap().Select(ex, q, start, end)
	assert.Nil(err)
	return &ex
}

// FetchSupplier select demand side
func FetchSupplier(start int64, end int64) *models.Exchange {
	ex := models.Exchange{}
	q := fmt.Sprintf(`SELECT SUM(request_in_count) AS supplier_impression_in,
	SUM(deliver_count) AS supplier_impression_out,
	SUM(deliver_bid) AS spent
	FROM %s
	WHERE time_id >= ?
	AND time_id <= ?`, models.SuplierTableName)
	m := models.NewManager()
	_, err := m.GetRDbMap().Select(ex, q, start, end)
	assert.Nil(err)
	return &ex
}

// ExchangeReport cron worker report exchange
func ExchangeReport(date time.Time) {
	//r,err := mysql.Manager.GetRDbMap().Select()
	start, end := factTableYesterdayID(date)
	dem := FetchDemand(start, end)
	sup := FetchSupplier(start, end)
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
				`, models.ExchangeTableName)
	m := models.NewManager()
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
	return datamodels.FactTableID(from), datamodels.FactTableID(to)
}
