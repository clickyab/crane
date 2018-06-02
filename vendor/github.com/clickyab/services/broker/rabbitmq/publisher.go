package rabbitmq

import (
	"sync"

	"github.com/clickyab/services/broker/rabbitmq/mq"
	"github.com/clickyab/services/broker/rabbitmq/mqfake"

	"github.com/clickyab/services/broker/rabbitmq/mqinterface"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/streadway/amqp"
)

type chnlLock struct {
	cnn    mqinterface.Connection
	chn    mqinterface.Channel
	lock   *sync.Mutex
	rtrn   chan amqp.Confirmation
	wg     *sync.WaitGroup
	closed bool
}

type consumer struct {
	amqp Amqp
}

// Publish try to publish an event
func (c consumer) Publish(in broker.Job) {
	rep := in.Report()
	var err error
	defer func() {
		rep(err)
	}()

	err = c.amqp.Publish(in)
	assert.Nil(err)
}

// NewRabbitBroker return a new rabbit broker
func NewRabbitBroker() broker.Interface {
	cns := consumer{
		amqp: Amqp{
			Channel:    &amqp.Channel{},
			Connection: &mq.Connection{},
			Dial:       &mq.Dial{},
		},
	}

	return &cns
}

// NewFakeRabbitBroker return a new rabbit broker
func NewFakeRabbitBroker() broker.Interface {
	cns := consumer{
		amqp: Amqp{
			Channel:    &mqfake.FakeChannel{},
			Connection: &mqfake.FakeConnection{},
			Dial:       &mqfake.FakeDial{},
		},
	}

	return &cns
}
