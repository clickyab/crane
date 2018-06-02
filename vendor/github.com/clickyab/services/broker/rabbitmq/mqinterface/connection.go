package mqinterface

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
}
