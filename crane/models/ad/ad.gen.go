// Code generated build with models DO NOT EDIT.

package ad

import (
	"github.com/clickyab/services/mysql"
	gorp "gopkg.in/gorp.v2"
)

// AUTO GENERATED CODE. DO NOT EDIT!

const (
	// AdTableFull is the Ad table name
	AdTableFull = "ads"
)

// Manager is the model manager for ad package
type Manager struct {
	mysql.Manager
}

// NewAdManager create and return a manager for this module
func NewAdManager() *Manager {
	return &Manager{}
}

// NewAdManagerFromTransaction create and return a manager for this module from a transaction
func NewAdManagerFromTransaction(tx gorp.SqlExecutor) (*Manager, error) {
	m := &Manager{}
	err := m.Hijack(tx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Initialize ad package
func (m *Manager) Initialize() {

	m.AddTableWithName(
		Ad{},
		AdTableFull,
	).SetKeys(
		true,
		"AdID",
	)

}
func init() {
	mysql.Register(NewAdManager())
}
