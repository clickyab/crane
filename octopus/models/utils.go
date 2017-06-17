package models

import (
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

func factTableRange(t time.Time) (int64, int64) {
	y, m, d := t.Date()
	from := time.Date(y, m, d, 0, 0, 1, 0, time.UTC)
	to := time.Date(y, m, d, 23, 59, 59, 0, time.UTC)
	return FactTableID(from), FactTableID(to)
}

func init() {
	layout := "2006-01-02T15:04:05.000Z"
	str := "2017-03-21T00:00:00.000Z"
	var err error
	epoch, err = time.Parse(layout, str)
	assert.Nil(err)
}
