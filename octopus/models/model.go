package models

import "time"

// SupplierReporter table
type SupplierReporter struct {
	ID                  int64     `json:"id" db:"id"`
	Supplier            string    `json:"supplier" db:"supplier"`
	Date                time.Time `json:"target_date" db:"target_date"`
	ImpressionIn        int64     `json:"impression_in" db:"impression_in"`
	ImpressionOut       int64     `json:"impression_out" db:"impression_out"`
	DeliveredImpression int64     `json:"delivered_impression" db:"delivered_impression"`
	Earn                int64     `json:"earn" db:"earn"`
}
