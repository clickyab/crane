package models

import (
	"clickyab.com/exchange/services/mysql"
)

// Manager is the model manager
type Manager struct {
	mysql.Manager
}

// Initialize the manager. nothing to do, just keep it in interface shape
func (m *Manager) Initialize() {
	m.AddTableWithName(
		SupplierSourceDemand{},
		"sup_dem_src",
	).SetKeys(
		true,
		"ID",
	)
	m.AddTableWithName(
		SupplierSource{},
		"sup_src",
	).SetKeys(
		true,
		"ID",
	)
	m.AddTableWithName(
		SupplierSource{},
		"exchange_report",
	).SetKeys(
		true,
		"ID",
	)

	m.AddTableWithName(DemandReport{}, "demand_report").
		SetKeys(true, "ID")

	m.AddTableWithName(
		SupplierReporter{},
		"supplier_report",
	).SetKeys(
		true,
		"ID")

}

// NewManager return a new manager object
func NewManager() *Manager {
	return &Manager{}
}

func init() {
	mysql.Register(NewManager())
}
