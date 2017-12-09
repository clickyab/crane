package config

import (
	"runtime"
	"sync"

	"github.com/clickyab/services/assert"

	onion "gopkg.in/fzerorubigd/onion.v3"
	_ "gopkg.in/fzerorubigd/onion.v3/yamlloader" // for loading yaml file
)

var (
	cfg appConfig
	// map of conf,description
	configs = make(map[string]string)
	locker  = &sync.Mutex{}
	o       = onion.New()
)

// AppConfig is the application config
type appConfig struct {
	DevelMode       bool   `onion:"devel_mode"`
	MaxCPUAvailable int    `onion:"max_cpu_available"`
	MachineName     string `onion:"machine_name"`
	// Setting time zone is not correct in app level. set it in system level
	//TimeZone        string `onion:"time_zone"`
}

// DescriptiveLayer is based on onion layer interface
type DescriptiveLayer interface {
	onion.Layer
	// Add get Description, key and value
	Add(string, string, interface{})
}

// NewDescriptiveLayer return new DescriptiveLayer
func NewDescriptiveLayer() DescriptiveLayer {
	return &layer{}
}

func defaultLayer() onion.Layer {
	d := onion.NewDefaultLayer()
	assert.Nil(d.SetDefault("core.devel_mode", true))
	assert.Nil(d.SetDefault("core.max_cpu_available", runtime.NumCPU()))
	//assert.Nil(d.SetDefault("core.time_zone", "Asia/Tehran"))
	assert.Nil(d.SetDefault("core.machine_name", "m1"))
	return d
}

// layer is configuration holder
type layer struct {
	defaltLayer onion.DefaultLayer
}

// Load a layer into the Onion. the call is only done in the
// registration
func (l *layer) Load() (map[string]interface{}, error) {
	if l.defaltLayer != nil {
		return l.defaltLayer.Load()
	}
	return map[string]interface{}{}, nil
}

// Add set a default value for a key
func (l *layer) Add(description string, key string, value interface{}) {
	locker.Lock()
	defer locker.Unlock()
	if l.defaltLayer == nil {
		l.defaltLayer = onion.NewDefaultLayer()
	}
	assert.Nil(l.defaltLayer.SetDefault(key, value))
	setDescription(key, description)

}

// GetDescriptions return config key, description
func GetDescriptions() map[string]string {
	return configs
}
