// Code generated build with models DO NOT EDIT.

package publisher

// AUTO GENERATED CODE. DO NOT EDIT!

import (
	"fmt"
	"strings"
	"time"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/initializer"
	gorp "gopkg.in/gorp.v2"
)

// CreatePublisher try to save a new Publisher in database
func (m *Manager) CreatePublisher(p *Publisher) error {
	now := time.Now()
	p.CreatedAt = &now
	p.UpdatedAt = &now
	func(in interface{}) {
		if ii, ok := in.(initializer.Simple); ok {
			ii.Initialize()
		}
	}(p)

	return m.GetWDbMap().Insert(p)
}

// UpdatePublisher try to update Publisher in database
func (m *Manager) UpdatePublisher(p *Publisher) error {
	now := time.Now()
	p.UpdatedAt = &now
	func(in interface{}) {
		if ii, ok := in.(initializer.Simple); ok {
			ii.Initialize()
		}
	}(p)

	_, err := m.GetWDbMap().Update(p)
	return err
}

// ListPublishersWithFilter try to list all Publishers without pagination
func (m *Manager) ListPublishersWithFilter(filter string, params ...interface{}) []Publisher {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	var res []Publisher
	_, err := m.GetRDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", PublisherTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListPublishers try to list all Publishers without pagination
func (m *Manager) ListPublishers() []Publisher {
	return m.ListPublishersWithFilter("")
}

// CountPublishersWithFilter count entity in Publishers table with valid where filter
func (m *Manager) CountPublishersWithFilter(filter string, params ...interface{}) int64 {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	cnt, err := m.GetRDbMap().SelectInt(
		fmt.Sprintf("SELECT COUNT(*) FROM %s %s", PublisherTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return cnt
}

// CountPublishers count entity in Publishers table
func (m *Manager) CountPublishers() int64 {
	return m.CountPublishersWithFilter("")
}

// ListPublishersWithPaginationFilter try to list all Publishers with pagination and filter
func (m *Manager) ListPublishersWithPaginationFilter(
	offset, perPage int, filter string, params ...interface{}) []Publisher {
	var res []Publisher
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}

	filter += " LIMIT ?, ? "
	params = append(params, offset, perPage)

	// TODO : better pagination without offset and limit
	_, err := m.GetRDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", PublisherTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListPublishersWithPagination try to list all Publishers with pagination
func (m *Manager) ListPublishersWithPagination(offset, perPage int) []Publisher {
	return m.ListPublishersWithPaginationFilter(offset, perPage, "")
}

// FindPublisherByFID return the Publisher base on its id
func (m *Manager) FindPublisherByFID(id int64) (*Publisher, error) {
	var res Publisher
	err := m.GetRDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE id=?", PublisherTableFull),
		id,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// PreInsert is gorp hook to prevent Insert without transaction
func (p *Publisher) PreInsert(q gorp.SqlExecutor) error {
	if _, ok := q.(*gorp.Transaction); !ok {
		// This is not good. gorp is not in transaction
		return fmt.Errorf("Insert Publisher must be in transaction")
	}
	return nil
}
