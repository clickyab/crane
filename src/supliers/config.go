package supliers

import (
	"services/config"

	"gopkg.in/fzerorubigd/onion.v2"
)

var (
	domain string
)

type cfgInitializer struct {
	o *onion.Onion
}

func (ci *cfgInitializer) Initialize(o *onion.Onion) []onion.Layer {
	ci.o = o
	l := onion.NewDefaultLayer()
	l.SetDefault("exchange.supplier.domain", "localhost")
	return []onion.Layer{l}
}

func (ci *cfgInitializer) Loaded() {
	domain = ci.o.GetStringDefault("exchange.supplier.domain", "localhost")
}

func init() {
	config.Register(&cfgInitializer{})
}
