package influx

import "github.com/clickyab/services/config"

var (
	server     = config.RegisterString("services.influx.server", "http://127.0.0.1:8600", "the influxdb server address")
	protocol   = config.RegisterString("services.influx.protocol", "http", "the influxdb server protocol, valids are http/udp")
	database   = config.RegisterString("services.influx.database", "clickyab", "the influxdb server database")
	user       = config.RegisterString("services.influx.user", "clickyab", "the influxdb server user")
	password   = config.RegisterString("services.influx.password", "bita123", "the influxdb server password")
	bufferSize = config.RegisterInt("services.influx.buffer_size", 1000, "buffer size for influxdb")
)
