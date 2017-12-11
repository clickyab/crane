package commands

import (
	"os"

	"github.com/clickyab/services/config"
)

const (
	// AppName the application name
	AppName string = "crane"
	// Organization the organization name
	Organization = "clickyab"
	// Prefix the prefix for config loader from env
	Prefix = "CRN"
)

// DefaultConfig for this set of apps
func DefaultConfig() config.DescriptiveLayer {
	d := config.NewDescriptiveLayer()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	d.Add("DESCRIPTION", "service.mysql.wdsn", "root:bita123@tcp(127.0.0.1:3306)/clickyab?charset=utf8&parseTime=true&charset=utf8")
	d.Add("DESCRIPTION", "service.mysql.rdsn", "root:bita123@tcp(127.0.0.1:3306)/clickyab?charset=utf8&parseTime=true&charset=utf8")
	d.Add("DESCRIPTION", "exchange.router.listen", ":"+port)
	return d
}
