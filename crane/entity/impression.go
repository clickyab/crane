package entity

import (
	"net"
	"net/http"
)

// ImpressionAttributes is the imp attr key
type ImpressionAttributes string

// Impression is the single impression object
type Impression interface {
	Request() *http.Request
	// TrackID return the random id of this imp object
	TrackID() string
	// ClientID is the key to identify client
	ClientID() string
	// IP return the client ip
	IP() net.IP
	// UserAgent return the client user agent
	UserAgent() string
	// Publisher return the publisher that this client is going into system from that
	Publisher() Publisher
	// Location of the request
	Location() Location
	// OS the os of requester if available
	OS() OS
	// Slots is the slot for this request
	Slots() []Slot
	// Category returns category obviously
	Category() []Category
	// Attributes return the impression specific attributes
	Attributes() map[string]interface{}
}
