package rabbitmq

import (
	"clickyab.com/exchange/services/broker"
	"clickyab.com/exchange/services/random"

	"sync/atomic"

	"clickyab.com/exchange/services/safe"

	"context"

	"time"

	"clickyab.com/exchange/services/assert"
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
	// TODO : limit on workers must match with this prefetch
	err = c.Qos(100, 0, false)
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
	safe.ContinuesGoRoutine(func(cnl context.CancelFunc) {
		consumerTag := <-random.ID
		delivery, err := c.Consume(q.Name, consumerTag, false, false, false, false, nil)
		if err != nil {
			cnl()
			assert.Nil(err) // I know its somehow redundant.
			return
		}
		logrus.Debug("Worker started")
		cn.consume(kill, cnl, consumer.Consume(), c, delivery, consumerTag)
	}, time.Second)
	return nil
}

func (consumer) consume(ctx context.Context, cnl context.CancelFunc, consumer chan<- broker.Delivery, c *amqp.Channel, delivery <-chan amqp.Delivery, consumerTag string) {
	atomic.SwapInt64(&hasConsumer, 1)
	done := ctx.Done()

	cErr := c.NotifyClose(make(chan *amqp.Error))
bigLoop:
	for {
		select {
		case job, ok := <-delivery:
			assert.True(ok, "[BUG] Channel is closed! why??")
			consumer <- &jsonDelivery{delivery: &job}
		case <-done:
			logrus.Debug("closing channel")
			// break the continues loop
			cnl()
			_ = c.Cancel(consumerTag, true)
			break bigLoop
		case e := <-cErr:
			logrus.Errorf("%T => %+v", *e, *e)
		}
	}
}
