package client

import (
	"services/config"

	onion "gopkg.in/fzerorubigd/onion.v2"
)

var (
	ip2lserver string
)

type cfgInitializer struct {
	o *onion.Onion
}

func (ci *cfgInitializer) Initialize(o *onion.Onion) []onion.Layer {
	ci.o = o
	l := onion.NewDefaultLayer()
	l.SetDefault("service.ip2location.client.endpoint", "127.0.0.1:8190")
	return []onion.Layer{l}
}

func (ci *cfgInitializer) Loaded() {
	ip2lserver = ci.o.GetStringDefault("service.ip2location.client.endpoint", "127.0.0.1:8190")
}

func init() {
	config.Register(&cfgInitializer{})
}
