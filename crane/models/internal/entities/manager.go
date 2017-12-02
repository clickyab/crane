package entities

import "github.com/clickyab/services/mysql"

type Manager struct {
	mysql.Manager
}

func NewManager() *Manager {
	return &Manager{}
}

// Initialize aaa package
func (m *Manager) Initialize() {

}

func init() {
	mysql.Register(NewManager())
}
