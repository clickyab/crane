package mysql

import (
	"database/sql"

	"github.com/clickyab/services/mysql"
	mm "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/sirupsen/logrus"
)

func newConnection(dsn string) (*sql.DB, error) {
	return sql.Open("mysql", dsn)
}

type logger struct {
	fields logrus.Fields
}

func (g logger) Print(v ...interface{}) {
	logrus.WithFields(g.fields).Debug(v)
}

func init() {
	mysql.RegisterConnectionFactory(newConnection)
	l := logger{
		fields: logrus.Fields{
			"marker": "mysql-driver",
			"type":   "err",
		},
	}
	mm.SetLogger(l)
}
