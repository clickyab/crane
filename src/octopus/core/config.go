package core

import (
	"services/assert"
	"services/config"
	"time"
)

var maximumTimeout time.Duration

type cfgInitializer struct {
}

func (ci *cfgInitializer) Initialize() config.DescriptiveLayer {
	l := config.NewDescriptiveLayer()
	l.Add("DESCRIPTION", "exchange.core.maximum_timeout", time.Second)
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
