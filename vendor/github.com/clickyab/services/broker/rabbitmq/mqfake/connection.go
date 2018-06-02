package mqfake

import (
	"github.com/clickyab/services/broker/rabbitmq/mqinterface"
	"github.com/sirupsen/logrus"
)

// FakeConnection is publisher connection
type FakeConnection struct {
	key string
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
