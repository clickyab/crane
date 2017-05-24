package aredis

import (
	"os"

	"clickyab.com/exchange/services/config"

	"regexp"

	"fmt"
)

var redisPattern = regexp.MustCompile("^redis://([^:]+):([^@]+)@([^:]+):([0-9]+)$")

var (
	network  *string
	address  *string
	password *string
	poolsize *int
	db       *int
)

func init() {
	var (
		port = "6379"
		host = "127.0.0.1"
		pass string
	)

	redisURL := os.Getenv("REDIS_URL")
	if all := redisPattern.FindStringSubmatch(redisURL); len(all) == 5 {
		port = all[4]
		host = all[3]
		pass = all[2]
	}

	network = config.RegisterString("services.redis.network", "tcp", "Redis network (normally tcp)")
	address = config.RegisterString("services.redis.address", fmt.Sprintf("%s:%s", host, port), "redis address host:port")
	password = config.RegisterString("services.redis.password", pass, "redis password")
	poolsize = config.RegisterInt("services.redis.poolsize", 200, "redis pool size")
	db = config.RegisterInt("services.redis.db", 1, "redis db number")
}
