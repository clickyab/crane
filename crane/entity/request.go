package entity

import (
	"net"
)

// Protocol is the scheme of http url
type Protocol string

const (
	// HTTP is scheme of request
	HTTP Protocol = "http"
	// HTTPS is scheme of request
	HTTPS Protocol = "https"
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
	Protocol() Protocol
	// UserAgent user agent
	UserAgent() string
	// Location return the location
	Location() Location
}
