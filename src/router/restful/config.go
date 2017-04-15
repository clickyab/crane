package restful

import (
	"services/config"

	"gopkg.in/fzerorubigd/onion.v2"
)

var (
	listenAddress string
	domain        string
)

type cfgInitializer struct {
	o *onion.Onion
}

func (ci *cfgInitializer) Initialize(o *onion.Onion) []onion.Layer {
	ci.o = o
	l := onion.NewDefaultLayer()
	l.SetDefault("exchange.router.listen", ":80")
	l.SetDefault("exchange.router.domain", "localhost")
	return []onion.Layer{l}
}

func (ci *cfgInitializer) Loaded() {
	listenAddress = ci.o.GetStringDefault("exchange.router.listen", ":80")
	domain = ci.o.GetStringDefault("exchange.router.domain", "localhost")
}

func init() {
	config.Register(&cfgInitializer{})
}
