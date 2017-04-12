package commands

import (
	"services/assert"

	"gopkg.in/fzerorubigd/onion.v2"
)

func DefaultConfig() onion.Layer {
	d := onion.NewDefaultLayer()
	assert.Nil(d.SetDefault("service.mysql.wdsn", "root:bita123@tcp(127.0.0.1:3306)/exchange?charset=utf8&parseTime=true&charset=utf8"))
	assert.Nil(d.SetDefault("service.mysql.rdsn", "root:bita123@tcp(127.0.0.1:3306)/exchange?charset=utf8&parseTime=true&charset=utf8"))
	assert.Nil(d.SetDefault("exchange.router.listen", ":8090"))

	return d
}
