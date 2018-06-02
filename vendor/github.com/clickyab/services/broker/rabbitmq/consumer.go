package rabbitmq

import (
	"github.com/clickyab/services/broker"

	"github.com/clickyab/services/config"
)

var prefetchCount = config.RegisterInt("services.broker.rabbitmq.prefetch", 100, "the prefetch count")

func (cn consumer) RegisterConsumer(consumer broker.Consumer) error {
	return cn.amqp.RegisterConsumer(consumer, prefetchCount.Int())
}
