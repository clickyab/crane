package supliers

import (
	"services/config"
)

var (
	domain string
)

type cfgInitializer struct {
}

func (ci *cfgInitializer) Initialize() config.DescriptiveLayer {
	l := config.NewDescriptiveLayer()
	l.Add("DESCRIPTION", "exchange.supplier.domain", "localhost")
	return l
}

func (ci *cfgInitializer) Loaded() {
	domain = config.GetStringDefault("exchange.supplier.domain", "localhost")
}

func init() {
	config.Register(&cfgInitializer{})
}
