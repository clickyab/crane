package mysql

import (
	"database/sql"

	"clickyab.com/exchange/services/assert"
	"clickyab.com/exchange/services/mysql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	// Mock is the mock object for use with mysql service
	Mock sqlmock.Sqlmock
)

func newConnection(string) (db *sql.DB, err error) {
	db, Mock, err = sqlmock.New()
	assert.Nil(err)

	return
}

func init() {
	mysql.RegisterConnectionFactory(newConnection)
}
