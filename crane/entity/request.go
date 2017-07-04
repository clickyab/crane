package entity

import (
	"net"
)

// Request get the request from request
type Request interface {
	// IP get the real ip from the request
	IP() net.IP
	// OS return os of the requested client
	OS() OS
	// Client get the client key (cop)
	ClientID() string
	// Protocol return http or https
	Protocol() string
	// UserAgent user agent
	UserAgent() string
	// Location return the location
	Location() Location
	// Attributes get attributes of different type requests
	Attributes() map[string]string
}
