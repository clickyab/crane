// Code generated build with models DO NOT EDIT.

package campaign

import (
	"github.com/clickyab/services/mysql"
	gorp "gopkg.in/gorp.v2"
)

// AUTO GENERATED CODE. DO NOT EDIT!

const (
	// campaignTableFull is the campaign table name
	campaignTableFull = "campaigns"
)

// Manager is the model manager for campaign package
type Manager struct {
	mysql.Manager
}

// NewCampaignManager create and return a manager for this module
func NewCampaignManager() *Manager {
	return &Manager{}
}

// NewCampaignManagerFromTransaction create and return a manager for this module from a transaction
func NewCampaignManagerFromTransaction(tx gorp.SqlExecutor) (*Manager, error) {
	m := &Manager{}
	err := m.Hijack(tx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Initialize campaign package
func (m *Manager) Initialize() {

	m.AddTableWithName(
		campaign{},
		campaignTableFull,
	).SetKeys(
		true,
		"CampainID",
	)

}
func init() {
	mysql.Register(NewCampaignManager())
}
