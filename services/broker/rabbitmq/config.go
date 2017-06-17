package rabbitmq

import (
	"os"

	"clickyab.com/exchange/services/config"
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
	dsn := os.Getenv("RABBITMQ_URL")
	if dsn == "" {
		dsn = "amqp://exchange:bita123@127.0.0.1:5672/cy"
	}
	d.Add("amqp dsn", "service.amqp.dsn", dsn)
	d.Add("amqp exchange to publish into", "service.amqp.exchange", "cy")
	d.Add("amqp publisher to publish into", "service.ampq.publisher", 30)
	d.Add("amqp confirm channel len", "service.amqp.confirm_len", 200)
	d.Add("amqp debug mode", "service.amqp.debug", false)
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
