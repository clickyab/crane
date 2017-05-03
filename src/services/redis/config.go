package aredis

import (
	"os"

	"services/config"

	"regexp"

	"fmt"
)

var cfg cfgLoader
var redisPattern = regexp.MustCompile("^redis://([^:]+):([^@]+)@([^:]+):([0-9]+)$")

type cfgLoader struct {
	Network  string `onion:"network"`
	Address  string `onion:"address"`
	Password string `onion:"password"`
	PoolSize int    `onion:"pool_size"`
	DB       int    `onion:"db"`
}

func (cl *cfgLoader) Initialize() config.DescriptiveLayer {

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
	l := config.NewDescriptiveLayer()

	l.Add("DESCRIPTION", "services.redis.network", "tcp")
	l.Add("DESCRIPTION", "services.redis.address", fmt.Sprintf("%s:%s", host, port))
	l.Add("DESCRIPTION", "services.redis.password", pass)
	l.Add("DESCRIPTION", "services.redis.poolsize", 200)
	l.Add("DESCRIPTION", "services.redis.db", 1)
	return l
}

func (cl *cfgLoader) Loaded() {
	config.GetStruct("services.redis", cl)
}

func init() {
	config.Register(&cfg)
}
