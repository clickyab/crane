package aredis

import (
	"services/assert"

	"services/config"

	onion "gopkg.in/fzerorubigd/onion.v2"
)

var cfg cfgLoader

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
	assert.Nil(d.SetDefault("service.redis.network", "tcp"))
	assert.Nil(d.SetDefault("service.redis.address", ":6379"))
	assert.Nil(d.SetDefault("service.redis.password", ""))
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
