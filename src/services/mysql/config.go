package mysql

import (
	"services/assert"

	"services/config"

	onion "gopkg.in/fzerorubigd/onion.v2"
)

var cfg cfgLoader

type cfgLoader struct {
	o *onion.Onion `onion:"-"`

	WDSN              string `onion:"wdsn"`
	RDSN              string `onion:"rdsn"`
	MaxConnection     int    `onion:"max_connection"`
	MaxIdleConnection int    `onion:"max_idle_connection"`

	DevelMode bool `onion:"-"`
}

func (cl *cfgLoader) Initialize(o *onion.Onion) []onion.Layer {
	cl.o = o

	d := onion.NewDefaultLayer()
	assert.Nil(d.SetDefault("service.mysql.wdsn", "root:bita123@tcp(127.0.0.1:3306)/clickyab?charset=utf8&parseTime=true&charset=utf8"))
	assert.Nil(d.SetDefault("service.mysql.rdsn", "root:bita123@tcp(127.0.0.1:3306)/clickyab?charset=utf8&parseTime=true&charset=utf8"))
	assert.Nil(d.SetDefault("service.mysql.max_connection", 30))
	assert.Nil(d.SetDefault("service.mysql.max_idle_connection", 5))

	return []onion.Layer{d}
}

func (cl *cfgLoader) Loaded() {
	cl.o.GetStruct("service.mysql", cl)
	cl.DevelMode = cl.o.GetBool("core.devel_mode")
}

func init() {
	config.Register(&cfg)
}
