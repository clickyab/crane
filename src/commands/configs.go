package commands

import (
	"services/config"
)

// DefaultConfig for this set of apps
func DefaultConfig() config.DescriptiveLayer {
	d := config.NewDescriptiveLayer()
	d.Add("DESCRIPTION", "service.mysql.wdsn", "root:bita123@tcp(127.0.0.1:3306)/exchange?charset=utf8&parseTime=true&charset=utf8")
	d.Add("DESCRIPTION", "service.mysql.rdsn", "root:bita123@tcp(127.0.0.1:3306)/exchange?charset=utf8&parseTime=true&charset=utf8")
	d.Add("DESCRIPTION", "exchange.router.listen", ":8090")
	return d
}
