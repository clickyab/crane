package core

import (
	"time"

	"clickyab.com/exchange/services/assert"
	"clickyab.com/exchange/services/config"
)

var maximumTimeout time.Duration

type cfgInitializer struct {
}

func (ci *cfgInitializer) Initialize() config.DescriptiveLayer {
	l := config.NewDescriptiveLayer()
	l.Add("maximum time to wait for demands to respond", "exchange.core.maximum_timeout", time.Second)
	return l
}

func (ci *cfgInitializer) Loaded() {
	maximumTimeout = config.GetDuration("exchange.core.maximum_timeout")
	assert.True(maximumTimeout > 0)
	assert.True(maximumTimeout < 10*time.Second)
}

func init() {
	config.Register(&cfgInitializer{})
}
