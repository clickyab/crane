package mqinterface

import "github.com/streadway/amqp"

// Connection is publisher connection
type Connection interface {
	/*
	   Channel opens a unique, concurrent server channel to process the bulk of AMQP
	   messages. Any error from methods on this receiver will render the receiver
	   invalid and a new Channel should be opened.

	*/
	Channel() (Channel, error)

	//Key return connection uniqeue key
	Key() string

	/*
		NotifyClose registers a listener for close events either initiated by an error
		accompaning a connection.close method or by a normal shutdown.

		On normal shutdowns, the chan will be closed.

		To reconnect after a transport or protocol error, register a listener here and
		re-run your setup process.

	*/
	NotifyClose(receiver chan *amqp.Error) chan *amqp.Error

	// IsClosed return bool status of connection
	IsClosed() bool

	//Closed set connection close state
	Closed()

	//SetOpen set connection open state
	SetOpen()
}
