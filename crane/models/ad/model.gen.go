// Code generated build with models DO NOT EDIT.

package ad

import (
	"fmt"
	"strings"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/initializer"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateAd try to save a new Ad in database
func (m *Manager) CreateAd(ad *Ad) error {

	func(in interface{}) {
		if ii, ok := in.(initializer.Simple); ok {
			ii.Initialize()
		}
	}(ad)

	return m.GetWDbMap().Insert(ad)
}

// UpdateAd try to update Ad in database
func (m *Manager) UpdateAd(ad *Ad) error {

	func(in interface{}) {
		if ii, ok := in.(initializer.Simple); ok {
			ii.Initialize()
		}
	}(ad)

	_, err := m.GetWDbMap().Update(ad)
	return err
}

// ListAdsWithFilter try to list all Ads without pagination
func (m *Manager) ListAdsWithFilter(filter string, params ...interface{}) []Ad {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	var res []Ad
	_, err := m.GetRDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", AdTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListAds try to list all Ads without pagination
func (m *Manager) ListAds() []Ad {
	return m.ListAdsWithFilter("")
}

// CountAdsWithFilter count entity in Ads table with valid where filter
func (m *Manager) CountAdsWithFilter(filter string, params ...interface{}) int64 {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	cnt, err := m.GetRDbMap().SelectInt(
		fmt.Sprintf("SELECT COUNT(*) FROM %s %s", AdTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return cnt
}

// CountAds count entity in Ads table
func (m *Manager) CountAds() int64 {
	return m.CountAdsWithFilter("")
}

// ListAdsWithPaginationFilter try to list all Ads with pagination and filter
func (m *Manager) ListAdsWithPaginationFilter(
	offset, perPage int, filter string, params ...interface{}) []Ad {
	var res []Ad
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}

	filter += " LIMIT ?, ? "
	params = append(params, offset, perPage)

	// TODO : better pagination without offset and limit
	_, err := m.GetRDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", AdTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListAdsWithPagination try to list all Ads with pagination
func (m *Manager) ListAdsWithPagination(offset, perPage int) []Ad {
	return m.ListAdsWithPaginationFilter(offset, perPage, "")
}

// FindAdByAdID return the Ad base on its id
func (m *Manager) FindAdByAdID(id int64) (*Ad, error) {
	var res Ad
	err := m.GetRDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE id=?", AdTableFull),
		id,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}
