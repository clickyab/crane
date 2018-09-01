package rabbitmq

import (
	"github.com/streadway/amqp"
	"sync"
	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
	)

type ccn struct {
	amqp *amqp.Connection
	sync.RWMutex
}
 func (c ccn) Connection () *amqp.Connection {
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
				amqp: c,
				RWMutex: sync.RWMutex{},
			}
			cnn.amqp.NotifyClose(errChn)

		}
		return err
	}, tryLimit.Duration())

	go func() {
		for {
			err := <-errChn
			cnn.Lock()

			logrus.Error(err)
			safe.Try(func() error {
				c, err := amqp.Dial(dsn.String())

				if err == nil {
					if e:= cnn.amqp.Close() ; e!=nil {
						logrus.Error(e)
					}
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
