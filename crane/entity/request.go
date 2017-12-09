package entity

import (
	"net"
)

// Protocol is the scheme of http url
type Protocol string

func (p Protocol) String() string {
	return string(p)
}

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
	// Protocol return http or https
	Protocol() Protocol
	// UserAgent user agent
	UserAgent() string
	// Location return the location
	Location() Location
	// IsMobile shows if its a phone
	IsMobile() bool
	// ISP returns request isp name
	ISP() string
	// EventPage is a string, generated only from multiple request (not one request to select multiple ad)
	EventPage() string
	// Is this request contain alexa?
	Alexa() bool
}
