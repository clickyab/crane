package models

import (
	"services/assert"
	"entity"
)

// Supplier is a supplier in our system
type Supplier struct {
	ID            int64  `json:"id" db:"id"`
	SName         string `json:"name" db:"name"`
	SType         string `json:"type" db:"type"`
	Key           string `json:"-" db:"key"`
	SFloorCPM     int64  `json:"floor_cpm" db:"floor_cpm"`
	SSoftFloorCPM int64  `json:"soft_floor_cpm" db:"soft_floor_cpm"`
	UnderFloor    int    `json:"under_floor" db:"under_floor"`
	Excluded      string `json:"-" db:"excluded"`
	Rate          int    `json:"-" db:"call_rate"`
}

func (s Supplier) Name() string {
	return s.SName
}

func (s Supplier) FloorCPM() int64 {
	return s.SFloorCPM
}

func (s Supplier) SoftFloorCPM() int64 {
	return s.SSoftFloorCPM
}

func (s Supplier) ExcludedDemands() []string {
	return nil
}

func (Supplier) CountryWhiteList() []entity.Country {
	return nil
}

func (s Supplier) CallRate() int {
	if s.Rate < 10 {
		return 10
	}
	if s.Rate > 100 {
		return 100
	}

	return s.Rate
}

func (s Supplier) Type() string {
	return s.SType
}

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
