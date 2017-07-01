// Code generated build with models DO NOT EDIT.

package publisher

import (
	"github.com/clickyab/services/mysql"
	gorp "gopkg.in/gorp.v2"
)

// AUTO GENERATED CODE. DO NOT EDIT!

const (
	// PublisherTableFull is the Publisher table name
	PublisherTableFull = "publishers"
)

// Manager is the model manager for publisher package
type Manager struct {
	mysql.Manager
}

// NewPublisherManager create and return a manager for this module
func NewPublisherManager() *Manager {
	return &Manager{}
}

// NewPublisherManagerFromTransaction create and return a manager for this module from a transaction
func NewPublisherManagerFromTransaction(tx gorp.SqlExecutor) (*Manager, error) {
	m := &Manager{}
	err := m.Hijack(tx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Initialize publisher package
func (m *Manager) Initialize() {

	m.AddTableWithName(
		Publisher{},
		PublisherTableFull,
	).SetKeys(
		true,
		"IDd",
	)

}
func init() {
	mysql.Register(NewPublisherManager())
}
