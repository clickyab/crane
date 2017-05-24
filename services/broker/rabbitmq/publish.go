package rabbitmq

import (
	"container/ring"
	"errors"
	"sync"

	"clickyab.com/exchange/services/random"

	"clickyab.com/exchange/services/broker"

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
}

// Publish try to publish an event
func (consumer) Publish(in broker.Job) {
	rep := in.Report()
	var err error
	defer func() {
		rep(err)
	}()
	rng = rng.Next()
	v := rng.Value.(*chnlLock)
	v.lock.Lock()
	defer v.lock.Unlock()
	if v.closed {
		err = errors.New("waiting for finalize, can not publish")
		return
	}

	msg, err := in.Encode()
	if err != nil {
		return
	}

	pub := amqp.Publishing{
		CorrelationId: <-random.ID,
		Body:          msg,
	}

	v.wg.Add(1)
	defer func() {
		// If the result is error, release the lock, there is no message to confirm!
		if err != nil {
			v.wg.Done()
		}
	}()
	topic := in.Topic()
	if cfg.Debug {
		topic = "debug." + topic
	}
	err = v.chn.Publish(cfg.Exchange, topic, true, false, pub)
}

// FinalizeWait is a function to wait for all publication to finish. after calling this,
// must not call the PublishEvent
func finalizeWait() {
	for i := 0; i < cfg.Publisher; i++ {
		rng = rng.Next()
		v := rng.Value.(*chnlLock)
		v.lock.Lock()
		// I know this lock release at the end, not after for, and this is ok
		defer v.lock.Unlock()

		v.closed = true
		v.wg.Wait()
		_ = v.chn.Close()
	}
}

func publishConfirm(cl *chnlLock) {
	for range cl.rtrn {
		cl.wg.Done()
	}
}

// NewRabbitBroker return a new rabbit broker
func NewRabbitBroker() broker.Interface {
	return &consumer{}
}
