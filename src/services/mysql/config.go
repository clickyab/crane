package mysql

import (
	"fmt"
	"os"
	"regexp"
	"services/config"
)

var mysqlPattern = regexp.MustCompile("^mysql://([^:]+):([^@]+)@([^:]+):([0-9]+)/([a-zA-Z0-9-_]+)$")
var (
	wdsn              *string
	rdsn              *string
	maxConnection     = config.RegisterInt("services.mysql.max_connection", 30, "max connection")
	maxIdleConnection = config.RegisterInt("services.mysql.max_idle_connection", 5, "max idle connection")

	develMode = config.RegisterBoolean("core.devel_mode", true, "development mode")
)

func init() {
	var (
		port     = "3306"
		host     = "127.0.0.1"
		database = "exchange"
		user     = "root"
		pass     = "bita123"
	)

	redisURL := os.Getenv("DATABASE_URL")
	if all := mysqlPattern.FindStringSubmatch(redisURL); len(all) == 6 {
		port = all[4]
		host = all[3]
		user = all[1]
		pass = all[2]
		database = all[5]
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", user, pass, host, port, database)
	wdsn = config.RegisterString("services.mysql.wdsn", dsn, "write database dsn")
	rdsn = config.RegisterString("services.mysql.wdsn", dsn, "read database dsn")
}
