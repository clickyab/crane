package mysql

import (
	"services/config"
)

var cfg cfgLoader

type cfgLoader struct {
	WDSN              string `onion:"wdsn"`
	RDSN              string `onion:"rdsn"`
	MaxConnection     int    `onion:"max_connection"`
	MaxIdleConnection int    `onion:"max_idle_connection"`

	DevelMode bool `onion:"-"`
}

func (cl *cfgLoader) Initialize() config.DescriptiveLayer {
	l := config.NewDescriptiveLayer()
	l.Add("Write database", "service.mysql.wdsn", "root:bita123@tcp(127.0.0.1:3306)/clickyab?charset=utf8&parseTime=true&charset=utf8")
	l.Add("Read database", "service.mysql.rdsn", "root:bita123@tcp(127.0.0.1:3306)/clickyab?charset=utf8&parseTime=true&charset=utf8")
	l.Add("Mysql maximum connectio", "service.mysql.max_connection", 30)
	l.Add("Mysql idle connection", "service.mysql.max_idle_connection", 5)
	return l
}

func (cl *cfgLoader) Loaded() {
	config.GetStruct("service.mysql", cl)
	cl.DevelMode = config.GetBool("core.devel_mode")
}

func init() {
	config.Register(&cfg)
}
