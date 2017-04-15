package models

import (
	"database/sql"
	"entity"
	"services/assert"
)

// Supplier is a supplier in our system
type Supplier struct {
	ID            int64          `json:"id" db:"id"`
	SName         string         `json:"name" db:"name"`
	SType         string         `json:"type" db:"type"`
	Key           string         `json:"-" db:"key"`
	SFloorCPM     int64          `json:"floor_cpm" db:"floor_cpm"`
	SSoftFloorCPM int64          `json:"soft_floor_cpm" db:"soft_floor_cpm"`
	UnderFloor    int            `json:"under_floor" db:"under_floor"`
	Excluded      sql.NullString `json:"-" db:"excluded"`
	Share         int            `json:"-" db:"share"`
}

// Name of this supplier
func (s Supplier) Name() string {
	return s.SName
}

// FloorCPM of this supplier
func (s Supplier) FloorCPM() int64 {
	return s.SFloorCPM
}

// SoftFloorCPM of this supplier
func (s Supplier) SoftFloorCPM() int64 {
	return s.SSoftFloorCPM
}

// ExcludedDemands of this supplire @TODO implement this
func (s Supplier) ExcludedDemands() []string {
	return nil
}

// CountryWhiteList is the country allowed by this supplier @TODO implement this
func (Supplier) CountryWhiteList() []entity.Country {
	return nil
}

// Type is the supplier type
func (s Supplier) Type() string {
	return s.SType
}

// GetSuppliers return all suppliers @TODO manage active/disable
func (m *Manager) GetSuppliers() map[string]Supplier {
	q := "SELECT * FROM suppliers"
	var res []Supplier
	_, err := m.GetRDbMap().Select(&res, q)
	assert.Nil(err)
	ret := make(map[string]Supplier, len(res))
	for i := range res {
		ret[res[i].Key] = res[i]
	}

	return ret
}
