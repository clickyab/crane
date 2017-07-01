// Code generated build with models DO NOT EDIT.

package user

import (
	"github.com/clickyab/services/mysql"
	gorp "gopkg.in/gorp.v2"
)

// AUTO GENERATED CODE. DO NOT EDIT!

const (
	// UserTableFull is the User table name
	UserTableFull = "users"
)

// Manager is the model manager for user package
type Manager struct {
	mysql.Manager
}

// NewUserManager create and return a manager for this module
func NewUserManager() *Manager {
	return &Manager{}
}

// NewUserManagerFromTransaction create and return a manager for this module from a transaction
func NewUserManagerFromTransaction(tx gorp.SqlExecutor) (*Manager, error) {
	m := &Manager{}
	err := m.Hijack(tx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Initialize user package
func (m *Manager) Initialize() {

	m.AddTableWithName(
		User{},
		UserTableFull,
	).SetKeys(
		true,
		"ID",
	)

}
func init() {
	mysql.Register(NewUserManager())
}
