package rabbitmq

import (
	"time"

	"github.com/clickyab/services/config"
)

var (
	dsn        = config.RegisterString("services.amqp.dsn", "amqp://server:bita123@127.0.0.1:5672/cy", "amqp dsn")
	exchange   = config.RegisterString("services.amqp.exchange", "cy", "amqp exchange to publish into")
	publisher  = config.RegisterInt("services.ampq.publisher", 30, "amqp publisher to publish into")
	confirmLen = config.RegisterInt("services.amqp.confirm_len", 200, "amqp confirm channel len")
	debug      = config.RegisterBoolean("services.amqp.debug", false, "amqp debug mode")
	tryLimit   = config.RegisterDuration("services.amqp.try_limit", time.Minute, "the limit to incremental try wait")
)
