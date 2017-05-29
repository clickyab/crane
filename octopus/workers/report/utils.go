package report

import (
	"fmt"
	"time"

	"clickyab.com/exchange/octopus/workers/internal"
	"clickyab.com/exchange/services/assert"
)

func factTableID(tm time.Time) int64 {
	layout := "2006-01-02T15:04:05.000Z"
	str := "2017-03-21T00:00:00.000Z"
	epoch, err := time.Parse(layout, str)
	assert.Nil(err)
	return int64(tm.Sub(epoch).Hours()) + 1
}

func translator(r internal.SupplierReporter) string {
	return fmt.Sprintf(`("%s","%s",%d,%d,%d,%d)`,
		r.Supplier,
		r.Date.Format("2006-01-02"),
		r.ImpressionOut,
		r.ImpressionIn,
		r.DeliveredImpression,
		r.Earn)
}
