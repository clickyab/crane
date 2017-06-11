package models

import (
	"time"

	"fmt"

	"clickyab.com/exchange/services/assert"
)

const (
	supSrcTable    string = `sup_src`
	supReportTable string = `supplier_report`

	insertQ = `INSERT INTO %s
	(supplier, date, impression_out, impression_in, delivered_impression, earn) VALUES %s
	ON DUPLICATE KEY UPDATE
	  impression_in=values(impression_in),
	  impression_out=values(impression_out),
	  earn=values(earn),
	  delivered_impression=values(delivered_impression)`

	getQ = `SELECT supplier,
      	sum(imp_in_count) as impression_in,
      	sum(imp_out_count) as impression_out,
      	sum(deliver_count) as delivered_impression,
      	sum(deliver_bid) as earn
		FROM %s
		where time_id BETWEEN %d AND %d
		GROUP BY supplier`
)

// SupplierReport is the func for chron
func SupplierReport(date time.Time) {
	id := FactTableID(date)

	reports := getData(id)

	parts := []string{}
	for i := range reports {
		reports[i].Date = date
		parts = append(parts, translator(reports[i]))
	}

	insertData(parts)

}

func getData(id int64) []SupplierReporter {
	getQuery := fmt.Sprintf(getQ, supSrcTable, id, id+23)

	reports := []SupplierReporter{}
	_, err := NewManager().GetRDbMap().Select(&reports, getQuery)
	assert.Nil(err)

	return reports
}

func insertData(parts []string) {

	m := NewManager()
	err := m.Begin()
	defer func() {
		if err != nil {
			assert.Nil(m.Rollback())
		} else {
			err = m.Commit()
		}
	}()
	if err != nil {
		return
	}

	for i := range parts {
		insertQuery := fmt.Sprintf(insertQ, supReportTable, parts[i])
		_, err := m.GetProperDBMap().Exec(insertQuery)
		if err != nil {
			return
		}
	}
}
