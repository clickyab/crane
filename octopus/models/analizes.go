package models

import "clickyab.com/exchange/services/assert"

// SupplierSourceDemand supplier_source_demand
type SupplierSourceDemand struct {
	ID         int64  `json:"id" db:"id"`
	Supplier   string `json:"supplier" db:"supplier"`
	Demand     string `json:"demand" db:"demand"`
	Source     string `json:"source" db:"source"`
	TimeID     int64  `json:"time_id" db:"time_id"`
	Request    int64  `json:"request_count" db:"request_count"`
	Impression int64  `json:"impression_count" db:"impression_count"`
	ShowTime   int64  `json:"show_count" db:"show_count"`
	ImpBid     int64  `json:"imp_bid" db:"imp_bid"`
	ShowBid    int64  `json:"show_bid" db:"show_bid"`
	WinBid     int64  `json:"win_bid" db:"win_bid"`
	Win        int64  `json:"win_count" db:"win_count"`
}

// SupplierSource supplier_source
type SupplierSource struct {
	ID         int64  `json:"id" db:"id"`
	Supplier   string `json:"supplier" db:"supplier"`
	Source     string `json:"source" db:"source"`
	TimeID     int64  `json:"time_id" db:"time_id"`
	Request    int64  `json:"request_count" db:"request_count"`
	Impression int64  `json:"impression_count" db:"impression_count"`
	ShowTime   int64  `json:"show_count" db:"show_count"`
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

// Parts is a multi query trick
type Parts struct {
	Query  string
	Params []interface{}
	Do     bool
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

	for i := range parts {
		if parts[i].Do {
			_, err = m.GetProperDBMap().Exec(parts[i].Query, parts[i].Params...)
			if err != nil {
				return
			}
		}
	}
	return nil
}
