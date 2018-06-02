package mqinterface

// Dial interface to implement dial to amqp
type Dial interface {
	// Dial accepts a string in the AMQP URI format and returns a new Connection
	// over TCP using PlainAuth.  Defaults to a server heartbeat interval of 10
	// seconds and sets the handshake deadline to 30 seconds. After handshake,
	// deadlines are cleared.
	//
	// Dial uses the zero value of tls.
	Dial(url string) (Connection, error)
}
