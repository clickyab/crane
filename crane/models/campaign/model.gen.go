// Code generated build with models DO NOT EDIT.

package campaign

// AUTO GENERATED CODE. DO NOT EDIT!

import (
	"fmt"
	"strings"
	"time"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/initializer"
	gorp "gopkg.in/gorp.v2"
)

// Createcampaign try to save a new campaign in database
func (m *Manager) Createcampaign(c *campaign) error {
	now := time.Now()
	c.CreatedAt = &now
	c.UpdatedAt = &now
	func(in interface{}) {
		if ii, ok := in.(initializer.Simple); ok {
			ii.Initialize()
		}
	}(c)

	return m.GetWDbMap().Insert(c)
}

// Updatecampaign try to update campaign in database
func (m *Manager) Updatecampaign(c *campaign) error {
	now := time.Now()
	c.UpdatedAt = &now
	func(in interface{}) {
		if ii, ok := in.(initializer.Simple); ok {
			ii.Initialize()
		}
	}(c)

	_, err := m.GetWDbMap().Update(c)
	return err
}

// ListcampaignsWithFilter try to list all campaigns without pagination
func (m *Manager) ListcampaignsWithFilter(filter string, params ...interface{}) []campaign {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	var res []campaign
	_, err := m.GetRDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", campaignTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// Listcampaigns try to list all campaigns without pagination
func (m *Manager) Listcampaigns() []campaign {
	return m.ListcampaignsWithFilter("")
}

// CountcampaignsWithFilter count entity in campaigns table with valid where filter
func (m *Manager) CountcampaignsWithFilter(filter string, params ...interface{}) int64 {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	cnt, err := m.GetRDbMap().SelectInt(
		fmt.Sprintf("SELECT COUNT(*) FROM %s %s", campaignTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return cnt
}

// Countcampaigns count entity in campaigns table
func (m *Manager) Countcampaigns() int64 {
	return m.CountcampaignsWithFilter("")
}

// ListcampaignsWithPaginationFilter try to list all campaigns with pagination and filter
func (m *Manager) ListcampaignsWithPaginationFilter(
	offset, perPage int, filter string, params ...interface{}) []campaign {
	var res []campaign
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}

	filter += " LIMIT ?, ? "
	params = append(params, offset, perPage)

	// TODO : better pagination without offset and limit
	_, err := m.GetRDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", campaignTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListcampaignsWithPagination try to list all campaigns with pagination
func (m *Manager) ListcampaignsWithPagination(offset, perPage int) []campaign {
	return m.ListcampaignsWithPaginationFilter(offset, perPage, "")
}

// FindcampaignByCampainID return the campaign base on its id
func (m *Manager) FindcampaignByCampainID(id int64) (*campaign, error) {
	var res campaign
	err := m.GetRDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE id=?", campaignTableFull),
		id,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// PreInsert is gorp hook to prevent Insert without transaction
func (c *campaign) PreInsert(q gorp.SqlExecutor) error {
	if _, ok := q.(*gorp.Transaction); !ok {
		// This is not good. gorp is not in transaction
		return fmt.Errorf("Insert campaign must be in transaction")
	}
	return nil
}
