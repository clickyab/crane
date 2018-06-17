package mqfake

import (
	"github.com/clickyab/services/broker/rabbitmq/mqinterface"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// FakeConnection is publisher connection
type FakeConnection struct {
	key    string
	closed bool
}

/*Channel opens a unique, concurrent server channel to process the bulk of AMQP
  messages. Any error from methods on this receiver will render the receiver
  invalid and a new Channel should be opened.

*/
func (f *FakeConnection) Channel() (mqinterface.Channel, error) {
	logrus.Debugln("create channel")

	return &FakeChannel{}, nil
}

//Key return connection uniqeue key
func (f *FakeConnection) Key() string {
	return f.key
}

/*NotifyClose registers a listener for close events either initiated by an error
accompaning a connection.close method or by a normal shutdown.

On normal shutdowns, the chan will be closed.

To reconnect after a transport or protocol error, register a listener here and
re-run your setup process.

*/
func (f *FakeConnection) NotifyClose(receiver chan *amqp.Error) chan *amqp.Error {
	return nil
}

// IsClosed return bool status of connection
func (f *FakeConnection) IsClosed() bool {
	return f.closed
}

//Closed set connection close state
func (f *FakeConnection) Closed() {
	f.closed = true
}

//SetOpen set connection open state
func (f *FakeConnection) SetOpen() {
	f.closed = false
}
