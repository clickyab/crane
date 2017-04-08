package models

import "services/mysql"

type Manager struct {
	mysql.Manager
}

func (m *Manager) Initialize() {

}

// NewManager return a new manager object
func NewManager() *Manager {
	return &Manager{}
}

func init() {
	mysql.Register(&Manager{})
}
