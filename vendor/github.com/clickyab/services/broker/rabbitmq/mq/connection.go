package mq

import (
	"sync"

	"github.com/clickyab/services/broker/rabbitmq/mqinterface"
	"github.com/streadway/amqp"
)

// Connection type to implement Connection interface
type Connection struct {
	amqpConnection amqp.Connection
	key            string
	url            string
	closed         bool

	lock sync.Mutex
}

/*Channel opens a unique, concurrent server channel to process the bulk of AMQP
  messages. Any error from methods on this receiver will render the receiver
  invalid and a new Channel should be opened.
*/
func (c *Connection) Channel() (mqinterface.Channel, error) {
	return c.amqpConnection.Channel()
}

//Key return connection uniqeue key
func (c *Connection) Key() string {
	return c.key
}

/*NotifyClose registers a listener for close events either initiated by an error
accompaning a connection.close method or by a normal shutdown.

On normal shutdowns, the chan will be closed.

To reconnect after a transport or protocol error, register a listener here and
re-run your setup process.

*/
func (c *Connection) NotifyClose(receiver chan *amqp.Error) chan *amqp.Error {
	return c.amqpConnection.NotifyClose(receiver)
}

//IsClosed return bool status of connection
func (c *Connection) IsClosed() bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.closed
}

//Closed set connection close state
func (c *Connection) Closed() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.closed = true
	return
}

//SetOpen set connection open state
func (c *Connection) SetOpen() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.closed = false
	return
}
