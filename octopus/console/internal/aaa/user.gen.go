// Code generated build with models DO NOT EDIT.

package aaa

// AUTO GENERATED CODE. DO NOT EDIT!

import (
	"fmt"
	"strings"

	"clickyab.com/exchange/services/assert"
	"clickyab.com/exchange/services/initializer"
	gorp "gopkg.in/gorp.v2"
)

// CreateUser try to save a new User in database
func (m *Manager) CreateUser(u *User) error {

	func(in interface{}) {
		if ii, ok := in.(initializer.Simple); ok {
			ii.Initialize()
		}
	}(u)

	return m.GetWDbMap().Insert(u)
}

// UpdateUser try to update User in database
func (m *Manager) UpdateUser(u *User) error {

	func(in interface{}) {
		if ii, ok := in.(initializer.Simple); ok {
			ii.Initialize()
		}
	}(u)

	_, err := m.GetWDbMap().Update(u)
	return err
}

// ListUsersWithFilter try to list all Users without pagination
func (m *Manager) ListUsersWithFilter(filter string, params ...interface{}) []User {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	var res []User
	_, err := m.GetRDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", UserTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListUsers try to list all Users without pagination
func (m *Manager) ListUsers() []User {
	return m.ListUsersWithFilter("")
}

// CountUsersWithFilter count entity in Users table with valid where filter
func (m *Manager) CountUsersWithFilter(filter string, params ...interface{}) int64 {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	cnt, err := m.GetRDbMap().SelectInt(
		fmt.Sprintf("SELECT COUNT(*) FROM %s %s", UserTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return cnt
}

// CountUsers count entity in Users table
func (m *Manager) CountUsers() int64 {
	return m.CountUsersWithFilter("")
}

// ListUsersWithPaginationFilter try to list all Users with pagination and filter
func (m *Manager) ListUsersWithPaginationFilter(
	offset, perPage int, filter string, params ...interface{}) []User {
	var res []User
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}

	filter += " LIMIT ?, ? "
	params = append(params, offset, perPage)

	// TODO : better pagination without offset and limit
	_, err := m.GetRDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", UserTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListUsersWithPagination try to list all Users with pagination
func (m *Manager) ListUsersWithPagination(offset, perPage int) []User {
	return m.ListUsersWithPaginationFilter(offset, perPage, "")
}

// FindUserByID return the User base on its id
func (m *Manager) FindUserByID(id int64) (*User, error) {
	var res User
	err := m.GetRDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE id=?", UserTableFull),
		id,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// FindUserByToken return the User base on its token
func (m *Manager) FindUserByToken(t string) (*User, error) {
	var res User
	err := m.GetRDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE token=?", UserTableFull),
		t,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// FindUserByEmail return the User base on its email
func (m *Manager) FindUserByEmail(e string) (*User, error) {
	var res User
	err := m.GetRDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE email=?", UserTableFull),
		e,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// PreInsert is gorp hook to prevent Insert without transaction
func (u *User) PreInsert(q gorp.SqlExecutor) error {
	if _, ok := q.(*gorp.Transaction); !ok {
		// This is not good. gorp is not in transaction
		return fmt.Errorf("Insert User must be in transaction")
	}
	return nil
}
