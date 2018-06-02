package mq

import (
	"github.com/clickyab/services/broker/rabbitmq/mqinterface"
	"github.com/streadway/amqp"
)

// Connection type to implement Connection interface
type Connection struct {
	amqpConnection amqp.Connection
	key            string
}

/*Channel opens a unique, concurrent server channel to process the bulk of AMQP
  messages. Any error from methods on this receiver will render the receiver
  invalid and a new Channel should be opened.
*/
func (c Connection) Channel() (mqinterface.Channel, error) {
	return c.amqpConnection.Channel()
}

//Key return connection uniqeue key
func (c Connection) Key() string {
	return c.key
}
