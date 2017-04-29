package aredis

import (
	"os"
	"services/assert"

	"services/config"

	"regexp"

	"fmt"

	onion "gopkg.in/fzerorubigd/onion.v2"
)

var cfg cfgLoader
var redisPattern = regexp.MustCompile("^redis://([^:]+):([^@]+)@([^:]+):([0-9]+)$")

type cfgLoader struct {
	o *onion.Onion `onion:"-"`

	Network  string `onion:"network"`
	Address  string `onion:"address"`
	Password string `onion:"password"`
	PoolSize int    `onion:"pool_size"`
	DB       int    `onion:"db"`
}

func (cl *cfgLoader) Initialize(o *onion.Onion) []onion.Layer {
	cl.o = o

	d := onion.NewDefaultLayer()
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

	assert.Nil(d.SetDefault("service.redis.network", "tcp"))
	assert.Nil(d.SetDefault("service.redis.address", fmt.Sprintf("%s:%s", host, port)))
	assert.Nil(d.SetDefault("service.redis.password", pass))
	assert.Nil(d.SetDefault("service.redis.poolsize", 200))
	assert.Nil(d.SetDefault("service.redis.db", 1))

	return []onion.Layer{d}
}

func (cl *cfgLoader) Loaded() {
	cl.o.GetStruct("service.redis", cl)
}

func init() {
	config.Register(&cfg)
}
