package config

import (
	"runtime"
	"services/assert"

	onion "gopkg.in/fzerorubigd/onion.v2"
)

var (
	o   = onion.New()
	cfg appConfig
)

// AppConfig is the application config
type appConfig struct {
	DevelMode       bool   `onion:"devel_mode"`
	MaxCPUAvailable int    `onion:"max_cpu_available"`
	MachineName     string `onion:"machine_name"`
	TimeZone        string `onion:"time_zone"`
}

func defaultLayer() onion.Layer {
	d := onion.NewDefaultLayer()
	assert.Nil(d.SetDefault("core.devel_mode", true))
	assert.Nil(d.SetDefault("core.max_cpu_available", runtime.NumCPU()))
	assert.Nil(d.SetDefault("core.time_zone", "Asia/Tehran"))
	assert.Nil(d.SetDefault("core.machine_name", "m1"))

	return d
}
