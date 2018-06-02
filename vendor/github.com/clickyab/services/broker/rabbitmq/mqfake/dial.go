package mqfake

import (
	"fmt"
	"time"

	"github.com/clickyab/services/broker/rabbitmq/mqinterface"
	"github.com/sirupsen/logrus"
)

// FakeDial interface to implement dial to amqp
type FakeDial struct {
}

// Dial accepts a string in the AMQP URI format and returns a new Connection
// over TCP using PlainAuth.  Defaults to a server heartbeat interval of 10
// seconds and sets the handshake deadline to 30 seconds. After handshake,
// deadlines are cleared.
//
// Dial uses the zero value of tls.
func (f *FakeDial) Dial(url string) (mqinterface.Connection, error) {
	logrus.Debugf("fake dial to %s", url)

	fcnn := FakeConnection{
		key: fmt.Sprintf("cn_%d", time.Now().UnixNano()),
	}

	return &fcnn, nil
}
