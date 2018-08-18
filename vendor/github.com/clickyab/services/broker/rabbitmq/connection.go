package rabbitmq

import (
	"github.com/streadway/amqp"
	"sync"
	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
	)

type ccn struct {
	amqp *amqp.Connection
	sync.Locker
}

func newConnection() *ccn {
	var cnn *ccn

	errChn := make(chan *amqp.Error)
	safe.Try(func() error {
		c, err := amqp.Dial(dsn.String())
		cnn.amqp.NotifyClose(errChn)
		if err == nil {
			cnn = &ccn{
				amqp: c,
			}
		}
		return err
	}, tryLimit.Duration())

	go func() {
		for {
			err := <-errChn
			logrus.Error(err)
			safe.Try(func() error {
				c, err := amqp.Dial(dsn.String())
				if err == nil {
					cnn.amqp = c
					c.NotifyClose(errChn)
				}
				return err
			}, tryLimit.Duration())
		}

	}()
	return cnn
}
