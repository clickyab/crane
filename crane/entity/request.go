package entity

import (
	"net"
	"time"
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
	// Timestamp return the time that this request arrived
	Timestamp() time.Time
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
	// EventPage is a string, generated only from multiple request (not one request to select multiple ad)
	EventPage() string
	// Is this request contain alexa?
	Alexa() bool
	// Referrer of the page
	Referrer() string
	// Parent is the page contain the ad
	Parent() string
	// Suspicious is means this request is Suspicious
	Suspicious() int
}
