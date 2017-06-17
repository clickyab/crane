package aredis

import (
	"os"

	"github.com/clickyab/services/config"

	"regexp"

	"fmt"

	"gopkg.in/fzerorubigd/onion.v3"
)

var redisPattern = regexp.MustCompile("^redis://([^:]+):([^@]+)@([^:]+):([0-9]+)$")

var (
	network  onion.String
	address  onion.String
	password onion.String
	poolsize onion.Int
	db       onion.Int
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
