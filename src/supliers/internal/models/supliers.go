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
	SShare        int            `json:"-" db:"share"`
}

// Name of this supplier
func (s Supplier) Name() string {
	return s.SName
}

// FloorCPM of this supplier
func (s Supplier) FloorCPM() int64 {
	i := s.SFloorCPM * int64(100+s.Share())
	return i / 100
}

// SoftFloorCPM of this supplier
func (s Supplier) SoftFloorCPM() int64 {
	i := s.SSoftFloorCPM * int64(100+s.Share())
	return i / 100
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

// Share of the supplier
func (s Supplier) Share() int {
	if s.SShare > 100 {
		return 100
	}
	if s.SShare < 1 {
		s.SShare = 1
	}
	return s.SShare
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
