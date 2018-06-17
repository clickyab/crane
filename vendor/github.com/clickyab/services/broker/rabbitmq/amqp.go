package rabbitmq

import (
	"container/ring"
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/clickyab/services/broker/rabbitmq/mqinterface"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/xlog"
	"github.com/manucorporat/try"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	retryMax = config.RegisterDuration("services.mysql.max_retry_connection", 5*time.Minute, "max time app should fallback to get mysql connection")

	connRng *ring.Ring
	pubRng  *ring.Ring

	publishCounter int64
	consumeCounter int64

	prLock = &sync.Mutex{}
	crLock = &sync.Mutex{}
)

// Amqp struct
type Amqp struct {
	Channel    mqinterface.Channel
	Connection mqinterface.Connection
	Dial       mqinterface.Dial
}

// MakeConnections Dial to amqp and make a connection pool
func (rb *Amqp) MakeConnections(count int) {
	connRng = ring.New(count)

	for index := 0; index < count; index++ {
		safe.Try(func() error {
			c, err := rb.Dial.Dial(dsn.String())
			if err == nil {
				connRng.Value = c
				automicConRngNext()
			}

			return err
		}, tryLimit.Duration())
	}
}

// ExchangeDeclare Declare new exchange. to be sure exchange is declared
func (rb *Amqp) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool) error {
	// conn := connRng.Value.(*amqp.Connection)
	automicConRngNext()
	chn, err := connRng.Value.(mqinterface.Connection).Channel()
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

// RegisterPublishers register new job publishers per connections
func (rb *Amqp) RegisterPublishers(count int) error {
	pubRng = ring.New(count)

	for index := 0; index < count; index++ {
		connection := connRng.Value.(mqinterface.Connection)
		pchn, err := connection.Channel()
		if err != nil {
			return err
		}

		rtrn := make(chan amqp.Confirmation, confirmLen.Int())
		err = pchn.Confirm(false)
		if err != nil {
			return err
		}

		pchn.NotifyPublish(rtrn)
		tmp := chnlLock{
			cnn:    connection,
			chn:    pchn,
			lock:   &sync.Mutex{},
			wg:     &sync.WaitGroup{},
			rtrn:   rtrn,
			closed: false,
		}
		go publishConfirm(&tmp)
		ctx := context.Background()
		safe.GoRoutine(ctx, func() {
			notifyChannelClose(ctx, &tmp)
		})

		pubRng.Value = &tmp

		automicConRngNext()
		automicPubRngNext()
	}

	return nil
}

// Publish job to amqp
func (rb *Amqp) Publish(in broker.Job) error {
	pubCh := pubRng.Value.(*chnlLock)
	pubCh.lock.Lock()

	defer pubCh.lock.Unlock()

	if pubCh.closed {
		cnt := pubRng.Len()
		for i := 0; i < cnt; i++ {
			automicPubRngNext()
			if !pubCh.closed {
				return rb.Publish(in)
			}
		}
		return errors.New("amqp channels is closed, can not publish")
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

	err = pubCh.chn.Publish(exchange.String(), topic, true, false, pub)
	if err == nil {
		atomic.AddInt64(&publishCounter, 1)
	} else {
		atomic.AddInt64(&publishCounter, -1)
	}

	return err
}

// FinalizeWait is a function to wait for all publication to finish. after calling this,
// must not call the PublishEvent
func (rb *Amqp) FinalizeWait() {
	for i := 0; i < pubRng.Len(); i++ {
		automicPubRngNext()
		v := pubRng.Value.(*chnlLock)
		v.lock.Lock()
		// I know this lock release at the end, not after for, and this is ok
		defer v.lock.Unlock()

		v.closed = true
		v.wg.Wait()
		_ = v.chn.Close()
	}
}

// RegisterConsumer register new cunsumer and fetch some of jobs
func (rb *Amqp) RegisterConsumer(consumer broker.Consumer, prefetchCount int) error {
	conn := connRng.Value.(mqinterface.Connection)
	c, err := conn.Channel()
	if err != nil {
		return err
	}

	err = rb.ExchangeDeclare(
		exchange.String(), // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
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
			atomic.AddInt64(&consumeCounter, -1)
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

func (rb *Amqp) consume(ctx context.Context, cnl context.CancelFunc, consumer chan<- broker.Delivery, c mqinterface.Channel, delivery <-chan amqp.Delivery, consumerTag string) {
	atomic.SwapInt64(&hasConsumer, 1)
	done := ctx.Done()

	cErr := c.NotifyClose(make(chan *amqp.Error))
bigLoop:
	for {
		select {
		case job, ok := <-delivery:
			assert.True(ok, "[BUG] Channel is closed! why??")
			atomic.AddInt64(&consumeCounter, 1)
			consumer <- &jsonDelivery{delivery: &job}
		case <-done:
			logrus.Debug("closing channel")
			// break the continues loop
			cnl()
			_ = c.Cancel(consumerTag, true)
			break bigLoop
		case e := <-cErr:
			atomic.AddInt64(&consumeCounter, -1)
			logrus.Errorf("%T => %+v", *e, *e)
		}
	}
}

func publishConfirm(cl *chnlLock) {
	for range cl.rtrn {
		cl.wg.Done()
	}
}

func notifyChannelClose(ctx context.Context, cl *chnlLock) {
	err := cl.chn.NotifyClose(make(chan *amqp.Error))
	er := <-err
	xlog.GetWithError(ctx, fmt.Errorf("amqp error code: %d, Reason: %s", er.Code, er.Reason)).Error("amqp channel closed")

	cl.closed = true
	if er.Code == 320 { //check if close connection notify fire after notify channel close
		cl.cnn.Closed()
	}

	safe.Try(func() error { return tryReconnect(cl) }, retryMax.Duration())
}

func tryReconnect(cl *chnlLock) error {
	connection := cl.cnn

	if connection.IsClosed() {
		return fmt.Errorf("channel connection is closed and we wait to reconnect")
	}

	var err error
	try.This(func() {
		logrus.Debug("try reconnect channel .... ")

		pchn, err := connection.Channel()
		if pchn == nil {
			err = fmt.Errorf("connection problem, we can not make channel")
		}

		if err == nil {
			rtrn := make(chan amqp.Confirmation, confirmLen.Int())
			err = pchn.Confirm(false)
			if err == nil {
				pchn.NotifyPublish(rtrn)

				cl.lock.Lock()
				defer cl.lock.Unlock()
				cl.chn = pchn
				cl.cnn = connection
				cl.rtrn = rtrn
				cl.closed = false

				logrus.Debug("ok channel reconnect")
			}
		}

	}).Catch(func(e try.E) {
		fmt.Println("catch error in try reconnect")
		err = fmt.Errorf("try error: %s", e)
	})

	return err
}

// ConnectionsCount get connections count
func (rb *Amqp) ConnectionsCount() int {
	return connRng.Len()
}

// PublishersCount get all publishers count
func (rb *Amqp) PublishersCount() int {
	return pubRng.Len()
}

// PublishersPerConnection get publishers count per connections
func (rb *Amqp) PublishersPerConnection() []int64 {
	pubPerCnn := make(map[string]int64, connRng.Len())

	for i := 0; i < pubRng.Len(); i++ {
		v := pubRng.Value.(*chnlLock)
		automicPubRngNext()

		pubPerCnn[v.cnn.Key()]++
	}

	var slc []int64
	for _, v := range pubPerCnn {
		slc = append(slc, v)
	}

	return slc
}

// JobStatiscs get all jobs statistics
func (rb *Amqp) JobStatiscs() map[string]int64 {
	st := map[string]int64{
		"Publish": atomic.LoadInt64(&publishCounter),
		"Consume": atomic.LoadInt64(&consumeCounter),
	}

	return st
}

func automicPubRngNext() {
	prLock.Lock()
	defer prLock.Unlock()
	pubRng = pubRng.Next()
}

func automicConRngNext() {
	crLock.Lock()
	defer crLock.Unlock()
	connRng = connRng.Next()
}
