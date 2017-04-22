package rabbit

import (
	"services/assert"

	"services/config"

	onion "gopkg.in/fzerorubigd/onion.v2"
)

var cfg cfgLoader

type cfgLoader struct {
	o *onion.Onion `onion:"-"`

	DSN        string `onion:"dsn"`
	Exchange   string `onion:"exchange"`
	Publisher  int    `onion:"publisher"`
	ConfirmLen int    `onion:"confirm_len"`
	Debug      bool   `onion:"debug"`
}

func (cl *cfgLoader) Initialize(o *onion.Onion) []onion.Layer {
	cl.o = o

	d := onion.NewDefaultLayer()
	assert.Nil(d.SetDefault("service.amqp.dsn", "amqp://server:bita123@127.0.0.1:5672/cy"))
	assert.Nil(d.SetDefault("service.amqp.exchange", "cy"))
	assert.Nil(d.SetDefault("service.ampq.publisher", 30))
	assert.Nil(d.SetDefault("service.amqp.confirm_len", 200))
	assert.Nil(d.SetDefault("service.amqp.debug", false))

	return []onion.Layer{d}
}

func (cl *cfgLoader) Loaded() {
	cl.o.GetStruct("service.amqp", cl)
	if cl.Publisher < 1 {
		cl.Publisher = 1
	}
}

func init() {
	config.Register(&cfg)
}
