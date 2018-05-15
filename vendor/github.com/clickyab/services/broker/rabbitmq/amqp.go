package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Amqp struct
type Amqp struct {
}

var (
	connectionPool = make(map[string]*amqp.Connection, 0)

	publisherPool = make(map[string]*chnlLock, 0)
)

// MakeConnections Dial to amqp and make a connection pool
func (rb *Amqp) MakeConnections(count int) ([]string, error) {
	var keys []string

	for i := 0; i < count; i++ {
		safe.Try(func() error {
			c, err := amqp.Dial(dsn.String())
			if err == nil {
				key := fmt.Sprintf("cn_%d", time.Now().UnixNano())
				keys = append(keys, key)
				connectionPool[key] = c
			}

			return err
		}, tryLimit.Duration())
	}

	return keys, nil
}

// ExchangeDeclare Declare new exchange. to be sure exchange is declared
func (rb *Amqp) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool) error {
	// conn := connRng.Value.(*amqp.Connection)
	chn, err := connectionPool[getOneKey(connectionPool)].Channel()
	if err != nil {
		return err
	}

	defer func() {
		assert.Nil(chn.Close())
	}()

	return chn.ExchangeDeclare(
		name,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		amqp.Table{},
	)
}

// MakePublisher register new job publisher
func (rb *Amqp) MakePublisher(connectionKey string) (string, error) {
	pchn, err := connectionPool[connectionKey].Channel()
	if err != nil {
		return "", err
	}

	rtrn := make(chan amqp.Confirmation, confirmLen.Int())
	err = pchn.Confirm(false)
	if err != nil {
		return "", err
	}

	pchn.NotifyPublish(rtrn)
	tmp := chnlLock{
		chn:    pchn,
		lock:   &sync.Mutex{},
		wg:     &sync.WaitGroup{},
		rtrn:   rtrn,
		closed: false,
	}
	go publishConfirm(&tmp)

	key := fmt.Sprintf("pub_%d", time.Now().UnixNano())
	publisherPool[key] = &tmp

	return key, nil
}

func getOneKey(data map[string]*amqp.Connection) string {
	for i := range data {
		return i
	}

	return ""
}

// Publish job to amqp
func (rb *Amqp) Publish(in broker.Job, pubKey string) error {
	pubCh := publisherPool[pubKey]

	pubCh.lock.Lock()
	defer pubCh.lock.Unlock()
	if pubCh.closed {
		return errors.New("waiting for finalize, can not publish")
	}

	msg, err := in.Encode()
	if err != nil {
		return err
	}

	pub := amqp.Publishing{
		CorrelationId: <-random.ID,
		Body:          msg,
	}

	pubCh.wg.Add(1)
	defer func() {
		// If the result is error, release the lock, there is no message to confirm!
		if err != nil {
			pubCh.wg.Done()
		}
	}()
	topic := in.Topic()
	if debug.Bool() {
		topic = "debug." + topic
	}

	return pubCh.chn.Publish(exchange.String(), topic, true, false, pub)
}

// FinalizeWait is a function to wait for all publication to finish. after calling this,
// must not call the PublishEvent
func (rb *Amqp) FinalizeWait() {
	for range publisherPool {
		rng = rng.Next()
		pubKey := rng.Value.(string)

		pubCh := publisherPool[pubKey]
		pubCh.lock.Lock()
		// I know this lock release at the end, not after for, and this is ok
		defer pubCh.lock.Unlock()

		pubCh.closed = true
		pubCh.wg.Wait()
		_ = pubCh.chn.Close()
	}
}

// RegisterConsumer register new cunsumer and fetch some of jobs
func (rb *Amqp) RegisterConsumer(consumer broker.Consumer, prefetchCount int) error {
	conn := connectionPool[getOneKey(connectionPool)]
	c, err := conn.Channel()
	if err != nil {
		return err
	}
	err = c.ExchangeDeclare(
		exchange.String(), // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)

	if err != nil {
		return err
	}
	qu := consumer.Queue()
	if debug.Bool() {
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
	err = c.Qos(prefetchCount, 0, false)
	if err != nil {
		return err
	}

	topic := consumer.Topic()
	if debug.Bool() {
		topic = "debug." + topic
	}
	err = c.QueueBind(
		q.Name,            // queue name
		topic,             // routing key
		exchange.String(), // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	safe.ContinuesGoRoutine(kill, func(cnl context.CancelFunc) time.Duration {
		consumerTag := <-random.ID
		delivery, err := c.Consume(q.Name, consumerTag, false, false, false, false, nil)
		if err != nil {
			cnl()
			assert.Nil(err) // I know its somehow redundant.
			return 0
		}
		logrus.Debug("Worker started")
		rb.consume(kill, cnl, consumer.Consume(kill), c, delivery, consumerTag)
		return time.Second
	})
	return nil
}

func (rb *Amqp) consume(ctx context.Context, cnl context.CancelFunc, consumer chan<- broker.Delivery, c *amqp.Channel, delivery <-chan amqp.Delivery, consumerTag string) {
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
