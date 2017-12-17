package config

import (
	"io/ioutil"
	"runtime"

	"github.com/clickyab/services/assert"
	"github.com/fzerorubigd/expand"
	"github.com/sirupsen/logrus"
	onion "gopkg.in/fzerorubigd/onion.v3"
	"gopkg.in/fzerorubigd/onion.v3/extraenv"
)

var (
	all []Initializer
)

// Initializer is the config initializer for module
type Initializer interface {
	// Initialize is called when the module is going to add its layer
	Initialize() DescriptiveLayer
	// Loaded inform the modules that all layer are ready
	Loaded()
}

//Initialize try to initialize config
func Initialize(organization, appName, prefix string, layers ...onion.Layer) {
	assert.Nil(o.AddLayer(defaultLayer()))

	for i := range all {
		nL := all[i].Initialize()
		_ = o.AddLayer(nL)
	}

	// now add the layer provided by app
	for i := range layers {
		_ = o.AddLayer(layers[i])
	}
	// Now load external config to overwrite them all.
	if err := o.AddLayer(onion.NewFileLayer("/etc/" + organization + "/" + appName + ".yaml")); err == nil {
		logrus.Infof("loading config from %s", "/etc/"+organization+"/"+appName+".yaml")
	}
	p, err := expand.Path("$HOME/." + organization + "/" + appName + ".yaml")
	if err == nil {
		if err = o.AddLayer(onion.NewFileLayer(p)); err == nil {
			logrus.Infof("loading config from %s", p)
		}
	}

	p, err = expand.Path("$PWD/configs/" + appName + ".yaml")
	if err == nil {
		if err = o.AddLayer(onion.NewFileLayer(p)); err == nil {
			logrus.Infof("loading config from %s", p)
		}
	}

	o.AddLazyLayer(extraenv.NewExtraEnvLayer(prefix))

	// load all registered variables
	o.Load()
	o.GetStruct("core", &cfg)
	// tell them that every thing is loaded
	for i := range all {
		all[i].Loaded()
	}
	SetConfigParameter()
}

// SetConfigParameter try to set the config parameter for the logrus base on config
func SetConfigParameter() {
	if cfg.DevelMode {
		// In development mode I need colors :) candy mode is GREAT!
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, DisableColors: false})
		logrus.SetLevel(logrus.DebugLevel)

	} else {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: false, DisableColors: true})
		logrus.SetLevel(logrus.WarnLevel)
		logrus.SetOutput(ioutil.Discard) // Discard the stdout logging
	}

	numcpu := cfg.MaxCPUAvailable
	if numcpu < 1 || numcpu > runtime.NumCPU() {
		numcpu = runtime.NumCPU()
	}

	runtime.GOMAXPROCS(numcpu)

	// Set global timezone
	//if l, err := time.LoadLocation(cfg.TimeZone); err == nil {
	//	time.Local = l
	//}
}

func setDescription(key, desc string) {
	lock.Lock()
	defer lock.Unlock()
	if d, ok := configs[key]; ok && d != "" && desc == "" {
		// if the new description is empty and the old one is not, ignore the new one
		return
	}
	configs[key] = desc
}

// Register a config module
func Register(i ...Initializer) {
	all = append(all, i...)
}
