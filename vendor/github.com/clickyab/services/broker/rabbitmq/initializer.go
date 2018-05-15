package rabbitmq

import (
	"container/ring"
	"context"
	"sync"
	"sync/atomic"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/initializer"

	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	connRng *ring.Ring
	once    = sync.Once{}
	kill    context.Context
)

var (
	hasConsumer int64

	connectionKeys          = make([]string, 0)
	publisherKeys           = make([]string, 0)
	publishersPerConnection = make(map[string][]string, 0)
)

type initRabbit struct {
	notifyCloser chan *amqp.Error
	amqp         RabbitInterface
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
		connRng = ring.New(cnt)

		var err error
		connectionKeys, err = in.amqp.MakeConnections(cnt)
		assert.Nil(err)

		for _, v := range connectionKeys {
			connRng.Value = v
			connRng = connRng.Next()
		}

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

		rng = ring.New(publisher.Int())
		for i := 0; i < publisher.Int(); i++ {
			connRng = connRng.Next()

			cncKey := connRng.Value.(string)
			pubKey, err := in.amqp.MakePublisher(cncKey)
			assert.Nil(err)

			publisherKeys = append(publisherKeys, pubKey)
			publishersPerConnection[cncKey] = append(publishersPerConnection[cncKey], pubKey)

			rng.Value = pubKey
			rng = rng.Next()
		}

		logrus.Debugln("Rabbit initialized")
		logrus.Debugln(in.Statistics())

		go func() {
			c := ctx.Done()
			assert.NotNil(c, "[BUG] context has no mean to cancel/deadline/timeout")
			<-c
			cnl() // the parent context normally take care of the children. but for idiotic linter :)
			in.finalize()
		}()
	})

}

func (in *initRabbit) Statistics() map[string]interface{} {
	st := map[string]interface{}{
		"connections": len(connectionKeys),
		"publishers":  len(publisherKeys),
		"publishers_per_connections": func() map[int]int {
			res := make(map[int]int, len(connectionKeys))
			for i, k := range connectionKeys {
				res[i] = len(publishersPerConnection[k])
			}

			return res
		}(),
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
	init := initRabbit{}
	init.amqp = &Amqp{}

	return &init
}

// NewFakeRabbitMQInitializer make a fake rabbit mq server
// select it via selector package
func NewFakeRabbitMQInitializer() initializer.Interface {
	init := initRabbit{}
	init.amqp = &FakeAmqp{}

	return &init
}
