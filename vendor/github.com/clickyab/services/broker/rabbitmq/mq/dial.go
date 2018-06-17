package mq

import (
	"context"
	"fmt"
	"time"

	"github.com/clickyab/services/xlog"
	"github.com/sirupsen/logrus"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/safe"

	"github.com/clickyab/services/broker/rabbitmq/mqinterface"
	"github.com/streadway/amqp"
)

var (
	retryMax = config.RegisterDuration("services.mysql.max_retry_connection", 5*time.Minute, "max time app should fallback to get mysql connection")
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
		return &Connection{}, err
	}

	newcn := Connection{
		amqpConnection: *cn,
		key:            fmt.Sprintf("cn_%d", time.Now().UnixNano()),
		url:            url,
		closed:         false,
	}

	ctx := context.Background()
	safe.GoRoutine(ctx, func() { notifyClose(ctx, cn, &newcn) })

	return &newcn, nil
}

func notifyClose(ctx context.Context, cn *amqp.Connection, ncn *Connection) {
	err := cn.NotifyClose(make(chan *amqp.Error))
	er := <-err
	xlog.GetWithError(ctx, fmt.Errorf("amqp error code: %d, Reason: %s", er.Code, er.Reason)).Error("amqp connection closed")
	ncn.Closed()

	safe.Try(func() error {
		logrus.Warn("try to reco connection $$$$ ")
		ncn.lock.Lock()
		defer ncn.lock.Unlock()

		aqcn, err := amqp.Dial(ncn.url)
		if err != nil {
			return err
		}

		ncn.amqpConnection = *aqcn
		ncn.SetOpen()

		return nil
	}, retryMax.Duration())
}
