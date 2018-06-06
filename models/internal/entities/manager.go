package entities

import "github.com/clickyab/services/mysql"

// Manager is database manager
type Manager struct {
	mysql.Manager
}

// NewManager create new manager
func NewManager() *Manager {
	return &Manager{}
}

// Initialize aaa package
func (m *Manager) Initialize() {
	m.AddTableWithName(
		Seat{},
		"seats",
	).SetKeys(
		true,
		"ID",
	)

	m.AddTableWithName(
		PublisherPage{},
		"publisher_pages",
	).SetKeys(
		true,
		"ID",
	)

	m.AddTableWithName(
		CreativesLocationsReport{},
		"creatives_locations_reports",
	).SetKeys(
		true,
		"ID",
	)
}

func init() {
	mysql.Register(NewManager())
}
