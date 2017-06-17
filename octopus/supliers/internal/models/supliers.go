package models

import (
	"clickyab.com/exchange/octopus/exchange"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/mysql"
)

// RendererFactory is a factory function for a supplier base on its type
type RendererFactory func(exchange.Supplier, string) exchange.Renderer

// Supplier is a supplier in our system
type Supplier struct {
	ID            int64                 `json:"id" db:"id"`
	SName         string                `json:"name" db:"name"`
	SType         string                `json:"type" db:"type"`
	Key           string                `json:"-" db:"key"`
	SFloorCPM     int64                 `json:"floor_cpm" db:"floor_cpm"`
	SSoftFloorCPM int64                 `json:"soft_floor_cpm" db:"soft_floor_cpm"`
	UnderFloor    int                   `json:"under_floor" db:"under_floor"`
	Excluded      mysql.StringJSONArray `json:"excluded_demands" db:"excluded_demands"`
	SShare        int                   `json:"-" db:"share"`
	SActive       int                   `json:"-" db:"active"`

	r exchange.Renderer
}

// Renderer return this supplier renderer
func (s Supplier) Renderer() exchange.Renderer {
	return s.r
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
	return s.Excluded
}

// CountryWhiteList is the country allowed by this supplier @TODO implement this
func (Supplier) CountryWhiteList() []exchange.Country {
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

// GetSuppliers return all suppliers
func (m *Manager) GetSuppliers(factory RendererFactory) map[string]Supplier {
	q := "SELECT * FROM suppliers WHERE active <> 0"
	var res []Supplier
	_, err := m.GetRDbMap().Select(&res, q)
	assert.Nil(err)
	ret := make(map[string]Supplier, len(res))
	for i := range res {
		res[i].r = factory(res[i], res[i].SType)
		ret[res[i].Key] = res[i]
	}

	return ret
}
