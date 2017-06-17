package commands

import (
	"os"

	"github.com/clickyab/services/config"
)

// DefaultConfig for this set of apps
func DefaultConfig() config.DescriptiveLayer {
	d := config.NewDescriptiveLayer()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	d.Add("DESCRIPTION", "service.mysql.wdsn", "root:bita123@tcp(127.0.0.1:3306)/exchange?charset=utf8&parseTime=true&charset=utf8")
	d.Add("DESCRIPTION", "service.mysql.rdsn", "root:bita123@tcp(127.0.0.1:3306)/exchange?charset=utf8&parseTime=true&charset=utf8")
	d.Add("DESCRIPTION", "exchange.router.listen", ":"+port)
	return d
}
