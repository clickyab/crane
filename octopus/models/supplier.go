package models

import (
	"time"

	"fmt"

	"clickyab.com/exchange/services/assert"
)

// updateSupplierReport will update supplier report (inclusive)
func updateSupplierReport(t time.Time) {
	td := t.Format("2006-01-02")
	from, to := factTableRange(t)
	var q = fmt.Sprintf(`INSERT INTO %s (
								supplier,
								target_date,
								impression_in,
								ad_out_count,
								delivered_count,
								earn
								)
							SELECT supplier,
							"%s",
							sum(imp_in_count),
							sum(ad_out_count),
							sum(deliver_count),
							sum(deliver_bid)
								FROM %s WHERE time_id BETWEEN %d AND %d
							GROUP BY supplier
							 ON DUPLICATE KEY UPDATE
							  supplier=VALUES(supplier),
							  target_date=VALUES(target_date),
							  impression_in_count=VALUES(impression_in_count),
							  ad_out_count=VALUES(ad_out_count),
							  delivered_count=VALUES(delivered_count),
							  earn=VALUES(earn)`, SupplierReportTableName, td, SupplierTableName, from, to)

	_, err := NewManager().GetRDbMap().Exec(q)
	assert.Nil(err)
}

// UpdateSupplierRange will update supplier report in range of two date (inclusive)
func UpdateSupplierRange(from time.Time, to time.Time) {
	if from.Unix() > to.Unix() {
		from, to = to, from
	}
	to = to.Add(24 * time.Hour)
	for from.Unix() < to.Unix() {
		updateSupplierReport(from)
		from = from.Add(time.Hour * 24)
	}
}
