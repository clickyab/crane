package mq

import (
	"fmt"
	"time"

	"github.com/clickyab/services/broker/rabbitmq/mqinterface"
	"github.com/streadway/amqp"
)

// Dial interface to implement dial to amqp
type Dial struct {
}

// Dial accepts a string in the AMQP URI format and returns a new Connection
// over TCP using PlainAuth.  Defaults to a server heartbeat interval of 10
// seconds and sets the handshake deadline to 30 seconds. After handshake,
// deadlines are cleared.
//
// Dial uses the zero value of tls.
func (f *Dial) Dial(url string) (mqinterface.Connection, error) {
	cn, err := amqp.Dial(url)
	if err != nil {
		return Connection{}, err
	}

	newcn := Connection{
		amqpConnection: *cn,
		key:            fmt.Sprintf("cn_%d", time.Now().UnixNano()),
	}

	return newcn, nil
}
