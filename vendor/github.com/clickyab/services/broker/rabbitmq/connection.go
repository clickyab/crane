package rabbitmq

import (
	"sync"

	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type ccn struct {
	amqp *amqp.Connection
	sync.RWMutex
}

func (c ccn) Connection() *amqp.Connection {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	return c.amqp
}

func newConnection() *ccn {
	var cnn *ccn

	errChn := make(chan *amqp.Error)
	safe.Try(func() error {
		c, err := amqp.Dial(dsn.String())
		if err == nil {
			cnn = &ccn{
				amqp:    c,
				RWMutex: sync.RWMutex{},
			}
			cnn.amqp.NotifyClose(errChn)

		}
		return err
	}, tryLimit.Duration())

	go func() {
		defer cnn.amqp.Close()
		for {
			err := <-errChn
			cnn.Lock()

			logrus.Errorf("rabbit new connection: ", err)
			if err := cnn.amqp.Close(); err != nil {
				logrus.Error(err)
			}

			safe.Try(func() error {
				c, err := amqp.Dial(dsn.String())
				if err == nil {
					cnn.amqp = c
					c.NotifyClose(errChn)
				}
				return err
			}, tryLimit.Duration())

			cnn.Unlock()

		}

	}()
	return cnn
}
