package models

import "services/mysql"

// Manager is the type to handle connection
type Manager struct {
	mysql.Manager
}

// Initialize is the base point to register tables if required
func (m *Manager) Initialize() {

}

// NewManager create and return a manager for this module
func NewManager() *Manager {
	return &Manager{}
}

func init() {
	mysql.Register(&Manager{})
}
