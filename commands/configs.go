package commands

import (
	"github.com/clickyab/services/config"
	// Activate slack
	_ "github.com/clickyab/services/slack"
)

const (
	// AppName the application name
	AppName = "crane"
	// Organization the organization name
	Organization = "clickyab"
	// Prefix the prefix for config loader from env
	Prefix = "CRN"
)

// DefaultConfig for this set of apps
func DefaultConfig() config.DescriptiveLayer {
	d := config.NewDescriptiveLayer()
	d.Add("", "services.broker.provider", "rabbitmq")
	d.Add("", "services.amqp.dsn", "amqp://crane:bita123@127.0.0.1:5672/")
	d.Add("", "services.amqp.exchange", "crane")
	d.Add("", "services.mysql.wdsn", "root:bita123@tcp(127.0.0.1:3306)/clickyab?charset=utf8&parseTime=true&charset=utf8")
	d.Add("", "services.mysql.rdsn", "root:bita123@tcp(127.0.0.1:3306)/clickyab?charset=utf8&parseTime=true&charset=utf8")
	return d
}
