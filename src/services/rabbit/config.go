package rabbit

import (
	"services/config"
)

var cfg cfgLoader

type cfgLoader struct {
	DSN        string `onion:"dsn"`
	Exchange   string `onion:"exchange"`
	Publisher  int    `onion:"publisher"`
	ConfirmLen int    `onion:"confirm_len"`
	Debug      bool   `onion:"debug"`
}

func (cl *cfgLoader) Initialize() config.DescriptiveLayer {
	d := config.NewDescriptiveLayer()
	d.Add("DESCRITION", "service.amqp.dsn", "amqp://server:bita123@127.0.0.1:5672/cy")
	d.Add("DESCRITION", "service.amqp.exchange", "cy")
	d.Add("DESCRITION", "service.ampq.publisher", 30)
	d.Add("DESCRITION", "service.amqp.confirm_len", 200)
	d.Add("DESCRITION", "service.amqp.debug", false)
	return d
}

func (cl *cfgLoader) Loaded() {
	config.GetStruct("service.amqp", cl)
	if cl.Publisher < 1 {
		cl.Publisher = 1
	}
}

func init() {
	config.Register(&cfg)
}
