package rabbitmq

import (
	"clickyab.com/exchange/services/broker"
	"clickyab.com/exchange/services/random"

	"sync/atomic"

	"clickyab.com/exchange/services/safe"

	"context"

	"github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

func (cn consumer) RegisterConsumer(consumer broker.Consumer) error {
	c, err := conn.Channel()
	if err != nil {
		return err
	}
	err = c.ExchangeDeclare(
		cfg.Exchange, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		return err
	}
	qu := consumer.Queue()
	if cfg.Debug {
		qu = "debug." + qu
	}
	q, err := c.QueueDeclare(qu, true, false, false, false, nil)
	if err != nil {
		return err
	}

	// prefetch count
	// **WARNING**
	// If ignore this, then there is a problem with rabbit. prefetch all jobs for this worker then.
	// the next worker get nothing at all!
	// **WARNING**
	err = c.Qos(1, 0, false)
	if err != nil {
		return err
	}

	topic := consumer.Topic()
	if cfg.Debug {
		topic = "debug." + topic
	}
	err = c.QueueBind(
		q.Name,       // queue name
		topic,        // routing key
		cfg.Exchange, // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}
	consumerTag := <-random.ID
	delivery, err := c.Consume(q.Name, consumerTag, false, false, false, false, nil)
	if err != nil {
		return err
	}
	logrus.Debug("Worker started")
	safe.GoRoutine(func() {
		cn.consume(kill, consumer.Consume(), c, delivery, consumerTag)
	})
	return nil
}

func (consumer) consume(ctx context.Context, consumer chan<- broker.Delivery, c *amqp.Channel, delivery <-chan amqp.Delivery, consumerTag string) {
	atomic.SwapInt64(&hasConsumer, 1)
	done := ctx.Done()
bigLoop:
	for {
		select {
		case job := <-delivery:
			d := jsonDelivery{delivery: &job}
			consumer <- d
		case <-done:
			_ = c.Cancel(consumerTag, true)
			break bigLoop
		}
	}
}
