package manager

import "clickyab.com/exchange/services/mysql"

// Manager is the manager for console
type Manager struct {
	mysql.Manager
}

// Initialize the model, its the interface, not really need this
func (m *Manager) Initialize() {
	m.AddTableWithName(
		User{},
		"user",
	).SetKeys(
		true,
		"ID",
	)
}

// NewManager return a new manager object
func NewManager() *Manager {
	return &Manager{}
}

func init() {
	mysql.Register(&Manager{})
}
