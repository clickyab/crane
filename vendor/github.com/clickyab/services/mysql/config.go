package mysql

import (
	"fmt"
	"os"
	"regexp"

	"time"

	"github.com/clickyab/services/config"
	"gopkg.in/fzerorubigd/onion.v3"
)

var pattern = regexp.MustCompile("^mysql://([^:]+):([^@]+)@([^:]+):([0-9]+)/([a-zA-Z0-9-_]+)$")

var (
	wdsn              onion.String
	rdsnSlice         onion.String
	maxConnection     = config.RegisterInt("services.mysql.max_connection", 30, "max connection")
	maxIdleConnection = config.RegisterInt("services.mysql.max_idle_connection", 5, "max idle connection")

	develMode = config.RegisterBoolean("core.devel_mode", true, "development mode")
	retryMax  = config.RegisterDuration("services.mysql.max_retry_connection", time.Minute, "max time app should fallback to get mysql connection")
	// CD is cool down, the time needed to sleep after each update
	rdbUpdateCD = config.RegisterDuration("services.mysql.max_retry_connection", time.Minute*2, "refresh read connection status after this amount of time")

	validSecondsSlaveBehind = config.RegisterInt64("services.mysql.replication.delay", 30, "max time slave can be behind")
	dbReplicated            = config.RegisterBoolean("services.mysql.replication.status", false, "is rdb replicated?")
	needWrite               = config.RegisterBoolean("services.mysql.need.write", true, "is this instance need write?")
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
	if all := pattern.FindStringSubmatch(redisURL); len(all) == 6 {
		port = all[4]
		host = all[3]
		user = all[1]
		pass = all[2]
		database = all[5]
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", user, pass, host, port, database)
	wdsn = config.RegisterString("services.mysql.wdsn", dsn, "write database dsn")
	rdsnSlice = config.RegisterString("services.mysql.rdsn", dsn, "comma separated read database dsn")
}
