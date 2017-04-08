package models

// Supplier is a supplier in our system
type Supplier struct {
	ID           int64  `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Type         string `json:"type" db:"type"`
	Key          string `json:"key" db:"key"`
	FloorCPM     int64  `json:"floor_cpm" db:"floor_cpm"`
	SoftFloorCPM int64  `json:"soft_floor_cpm" db:"soft_floor_cpm"`
	UnderFloor   int    `json:"under_floor" db:"under_floor"`
}
