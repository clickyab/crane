package rabbitmq

import (
	"container/ring"
	"fmt"
	"sync"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/streadway/amqp"
)

var (
	rng *ring.Ring
)

type chnlLock struct {
	chn    Channel
	lock   *sync.Mutex
	rtrn   chan amqp.Confirmation
	wg     *sync.WaitGroup
	closed bool
}

type consumer struct {
	amqp RabbitInterface
}

// Publish try to publish an event
func (c consumer) Publish(in broker.Job) {
	rep := in.Report()
	var err error
	defer func() {
		rep(err)
	}()
	rng = rng.Next()
	pubKey := rng.Value.(string)

	fmt.Println("____________________________________________")
	fmt.Println("in publish")
	err = c.amqp.Publish(in, pubKey)
	fmt.Println("____________________________________________")
	assert.Nil(err)
}

func publishConfirm(cl *chnlLock) {
	for range cl.rtrn {
		cl.wg.Done()
	}
}

// NewRabbitBroker return a new rabbit broker
func NewRabbitBroker() broker.Interface {
	cns := consumer{}
	cns.amqp = &Amqp{}

	return &cns
}

// NewFakeRabbitBroker return a new rabbit broker
func NewFakeRabbitBroker() broker.Interface {
	cns := consumer{}
	cns.amqp = &FakeAmqp{}

	return &cns
}
