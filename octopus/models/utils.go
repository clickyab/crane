package models

import (
	"fmt"
	"time"

	"clickyab.com/exchange/services/assert"
)

var (
	epoch time.Time
)

// FactTableID is a helper function to get the fact table id from time
func FactTableID(tm time.Time) int64 {
	return int64(tm.Sub(epoch).Hours()) + 1
}

func translator(r SupplierReporter) string {
	return fmt.Sprintf(`("%s","%s",%d,%d,%d,%d)`,
		r.Supplier,
		r.Date.Format("2006-01-02"),
		r.ImpressionOut,
		r.ImpressionIn,
		r.DeliveredImpression,
		r.Earn)
}

func init() {
	layout := "2006-01-02T15:04:05.000Z"
	str := "2017-03-21T00:00:00.000Z"
	var err error
	epoch, err = time.Parse(layout, str)
	assert.Nil(err)
}
