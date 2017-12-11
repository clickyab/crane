package rabbitmq

import (
	"container/ring"
	"context"
	"sync"
	"sync/atomic"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/initializer"

	"fmt"

	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	conn *amqp.Connection
	once = sync.Once{}
	kill context.Context
)

// Channel opens a unique, concurrent server channel to process the bulk of AMQP
// messages. Any error from methods on this receiver will render the receiver
// invalid and a new Channel should be opened.
type Channel interface {
	/*
		Confirm puts this channel into confirm mode so that the client can ensure all
		publishing's have successfully been received by the server. After entering this
		mode, the server will send a basic.ack or basic.nack message with the deliver
		tag set to a 1 based incrementing index corresponding to every publishing
		received after the this method returns.

		Add a listener to Channel.NotifyPublish to respond to the Confirmations. If
		Channel.NotifyPublish is not called, the Confirmations will be silently
		ignored.

		The order of acknowledgments is not bound to the order of deliveries.

		Ack and Nack confirmations will arrive at some point in the future.

		Unroutable mandatory or immediate messages are acknowledged immediately after
		any Channel.NotifyReturn listeners have been notified. Other messages are
		acknowledged when all queues that should have the message routed to them have
		either have received acknowledgment of delivery or have enqueued the message,
		persisting the message if necessary.

		When noWait is true, the client will not wait for a response. A channel
		exception could occur if the server does not support this method.

	*/
	Confirm(noWait bool) error

	/*
	   NotifyPublish registers a listener for reliable publishing. Receives from this
	   chan for every publish after Channel.Confirm will be in order starting with
	   DeliveryTag 1.

	   There will be one and only one Confirmation Publishing starting with the
	   delivery tag of 1 and progressing sequentially until the total number of
	   publishing's have been seen by the server.

	   Acknowledgments will be received in the order of delivery from the
	   NotifyPublish channels even if the server acknowledges them out of order.

	   The listener chan will be closed when the Channel is closed.

	   The capacity of the chan Confirmation must be at least as large as the
	   number of outstanding publishing's. Not having enough buffered chans will
	   create a deadlock if you attempt to perform other operations on the Connection
	   or Channel while confirms are in-flight.

	   It's advisable to wait for all Confirmations to arrive before calling
	   Channel.Close() or Connection.Close().

	*/
	NotifyPublish(confirm chan amqp.Confirmation) chan amqp.Confirmation

	/*
	   Publish sends a Publishing from the client to an exchange on the server.

	   When you want a single message to be delivered to a single queue, you can
	   publish to the default exchange with the routingKey of the queue name. This is
	   because every declared queue gets an implicit route to the default exchange.

	   Since publishing's are asynchronous, any undeliverable message will get returned
	   by the server. Add a listener with Channel.NotifyReturn to handle any
	   undeliverable message when calling publish with either the mandatory or
	   immediate parameters as true.

	   publishing's can be undeliverable when the mandatory flag is true and no queue is
	   bound that matches the routing key, or when the immediate flag is true and no
	   consumer on the matched queue is ready to accept the delivery.

	   This can return an error when the channel, connection or socket is closed. The
	   error or lack of an error does not indicate whether the server has received this
	   publishing.

	   It is possible for publishing to not reach the broker if the underlying socket
	   is shutdown without pending publishing packets being flushed from the kernel
	   buffers. The easy way of making it probable that all publishing's reach the
	   server is to always call Connection.Close before terminating your publishing
	   application. The way to ensure that all publishing's reach the server is to add
	   a listener to Channel.NotifyPublish and put the channel in confirm mode with
	   Channel.Confirm. Publishing delivery tags and their corresponding
	   confirmations start at 1. Exit when all publishing's are confirmed.

	   When Publish does not return an error and the channel is in confirm mode, the
	   internal counter for DeliveryTags with the first confirmation starting at 1.

	*/
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error

	/*
	   Close initiate a clean channel closure by sending a close message with the error
	   code set to '200'.

	   It is safe to call this method multiple times.

	*/
	Close() error
}

// Connection is publisher connection
type Connection interface {
	/*
	   Channel opens a unique, concurrent server channel to process the bulk of AMQP
	   messages. Any error from methods on this receiver will render the receiver
	   invalid and a new Channel should be opened.

	*/
	Channel() (Channel, error)
}

var (
	hasConsumer int64
)

type initRabbit struct {
	notifyCloser chan *amqp.Error
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
		kill, _ = context.WithCancel(ctx)
		safe.Try(func() error {
			var err error
			conn, err = amqp.Dial(dsn.String())
			return err
		}, tryLimit.Duration())

		chn, err := conn.Channel()
		assert.Nil(err)
		defer chn.Close()

		assert.Nil(
			chn.ExchangeDeclare(
				exchange.String(),
				"topic",
				true,
				false,
				false,
				false,
				amqp.Table{},
			),
		)

		rng = ring.New(publisher.Int())
		for i := 0; i < publisher.Int(); i++ {
			pchn, err := conn.Channel()
			assert.Nil(err)
			rtrn := make(chan amqp.Confirmation, confirmLen.Int())
			err = pchn.Confirm(false)
			assert.Nil(err)
			pchn.NotifyPublish(rtrn)
			tmp := chnlLock{
				chn:    pchn,
				lock:   &sync.Mutex{},
				wg:     &sync.WaitGroup{},
				rtrn:   rtrn,
				closed: false,
			}
			go publishConfirm(&tmp)
			rng.Value = &tmp
			rng = rng.Next()
		}

		logrus.Debug("Rabbit initialized")

		go func() {
			c := ctx.Done()
			assert.NotNil(c, "[BUG] context has no mean to cancel/deadline/timeout")
			<-c
			finalize()
		}()
	})

}

// finalize try to close rabbitmq connection
func finalize() {
	if atomic.CompareAndSwapInt64(&hasConsumer, 1, 0) {

	}
	finalizeWait()
	logrus.Debug("Rabbit finalized.")
}

func init() {
	initializer.Register(&initRabbit{}, 0)
}
