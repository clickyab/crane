package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/clickyab/services/broker/rabbitmq/mqfake"
	"github.com/clickyab/services/healthz"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker/rabbitmq/mq"
	"github.com/clickyab/services/initializer"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	once        = sync.Once{}
	kill        context.Context
	hasConsumer int64
)

type initRabbit struct {
	notifyCloser chan *amqp.Error
	amqp         *Amqp
}

func (in *initRabbit) Healthy(context.Context) error {
	select {
	case err := <-in.notifyCloser:
		if err != nil {
			return fmt.Errorf("RabbitMQ error happen : %s", err)
		}
	default: // Do not block
	}
	return nil
}

// Initialize the module at the beginning of the application to create a publish channel
func (in *initRabbit) Initialize(ctx context.Context) {

	once.Do(func() {
		// the size is here for channel to not block the caller. since we read this on the health check command
		in.notifyCloser = make(chan *amqp.Error, 10)
		var cnl context.CancelFunc
		kill, cnl = context.WithCancel(ctx)
		cnt := connection.Int()
		if cnt < 1 {
			cnt = 1
		}

		in.amqp.MakeConnections(cnt)

		assert.Nil(
			in.amqp.ExchangeDeclare(
				exchange.String(),
				"topic",
				true,
				false,
				false,
				false,
			),
		)

		err := in.amqp.RegisterPublishers(publisher.Int())
		assert.Nil(err)

		logrus.Debugln("Rabbit initialized")
		logrus.Debugln(in.Statistics())

		go func() {
			c := ctx.Done()
			assert.NotNil(c, "[BUG] context has no mean to cancel/deadline/timeout")
			<-c
			cnl() // the parent context normally take care of the children. but for idiotic linter :)
			in.finalize()
		}()

		healthz.Register(in)
	})
}

func (in *initRabbit) Statistics() map[string]interface{} {
	st := map[string]interface{}{
		"connections": in.amqp.ConnectionsCount(),
		"publishers":  in.amqp.PublishersCount(),
		"jobs":        in.amqp.JobStatiscs(),
		"publishers_per_connections": in.amqp.PublishersPerConnection(),
	}

	return st
}

// finalize try to close rabbitmq connection
func (in *initRabbit) finalize() {
	if atomic.CompareAndSwapInt64(&hasConsumer, 1, 0) {

	}
	in.amqp.FinalizeWait()
	logrus.Debug("Rabbit finalized.")
}

// NewRabbitMQInitializer rabbit mq is an exception in domino. we need to register it whenever we
// select it via selector package
func NewRabbitMQInitializer() initializer.Interface {
	init := initRabbit{
		amqp: &Amqp{
			Channel:    &amqp.Channel{},
			Connection: &mq.Connection{},
			Dial:       &mq.Dial{},
		},
	}

	return &init
}

// NewFakeRabbitMQInitializer make a fake rabbit mq server
// select it via selector package
func NewFakeRabbitMQInitializer() initializer.Interface {
	init := initRabbit{
		amqp: &Amqp{
			Channel:    &mqfake.FakeChannel{},
			Connection: &mqfake.FakeConnection{},
			Dial:       &mqfake.FakeDial{},
		},
	}

	return &init
}
