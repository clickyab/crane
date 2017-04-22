package core

import (
	"services/assert"
	"services/config"
	"time"

	"gopkg.in/fzerorubigd/onion.v2"
)

var maximumTimeout time.Duration

type cfgInitializer struct {
	o *onion.Onion
}

func (ci *cfgInitializer) Initialize(o *onion.Onion) []onion.Layer {
	ci.o = o
	l := onion.NewDefaultLayer()
	l.SetDefault("exchange.core.maximum_timeout", time.Second)
	return []onion.Layer{l}
}

func (ci *cfgInitializer) Loaded() {
	maximumTimeout = ci.o.GetDuration("exchange.core.maximum_timeout")
	assert.True(maximumTimeout > 0)
	assert.True(maximumTimeout < 10*time.Second)
}

func init() {
	config.Register(&cfgInitializer{})
}
