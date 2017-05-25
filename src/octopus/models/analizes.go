package models

import "services/assert"

// SupplierSourceDemand supplier_source_demand
type SupplierSourceDemand struct {
	Supplier   string `json:"supplier" db:"supplier"`
	Demand     string `json:"demand" db:"demand"`
	Source     string `json:"source" db:"source"`
	TimeID     int64  `json:"time_id" db:"time_id"`
	Request    int64  `json:"request" db:"request"`
	Impression int64  `json:"impression" db:"impression"`
	ShowTime   int64  `json:"show_time" db:"show_time"`
	ImpBid     int64  `json:"imp_bid" db:"imp_bid"`
	ShowBid    int64  `json:"show_bid" db:"show_bid"`
	WinBid     int64  `json:"win_bid" db:"win_bid"`
	Win        int64  `json:"win" db:"win"`
}

// SupplierSource supplier_source
type SupplierSource struct {
	Supplier   string `json:"supplier" db:"supplier"`
	Source     string `json:"source" db:"source"`
	TimeID     int64  `json:"time_id" db:"time_id"`
	Request    int64  `json:"request" db:"request"`
	Impression int64  `json:"impression" db:"impression"`
	ShowTime   int64  `json:"show_time" db:"show_time"`
	ImpBid     int64  `json:"imp_bid" db:"imp_bid"`
	ShowBid    int64  `json:"show_bid" db:"show_bid"`
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

// UpdateConsume try to update the proper report tables
func (m *Manager) UpdateConsume(q1, q2 string) (err error) {
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
	_, err = m.GetProperDBMap().Exec(q1)
	if err != nil {
		return
	}
	_, err = m.GetProperDBMap().Exec(q2)
	if err != nil {
		return
	}
	return
}
