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

}

func init() {
	mysql.Register(NewManager())
}
