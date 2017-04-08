package entity

import (
	"net"
	"net/http"
)

// ImpressionType is the publisher type
type ImpressionType string

// ImpressionAttributes is the imp attr key
type ImpressionAttributes string

const (
	// ImpressionTypeApp is the app
	ImpressionTypeApp ImpressionType = "app"
	// ImpressionTypeWeb is the web
	ImpressionTypeWeb ImpressionType = "web"
	// ImpressionTypeVast is the vast
	ImpressionTypeVast ImpressionType = "vast"
)

// Impression is the single impression object
type Impression interface {
	Request() *http.Request
	// MegaIMP return the random id of this imp object
	MegaIMP() string
	// ClientID is the key to identify client
	ClientID() int64
	// IP return the client ip
	IP() net.IP
	// UserAgent return the client user agent
	UserAgent() string
	// Source return the publisher that this client is going into system from that
	Source() Publisher
	// Location of the request
	Location() Location
	// OS the os of requester if available
	OS() OS
	// Attributes is the generic attribute system
	Attributes(ImpressionAttributes) interface{}
	// Slots is the slot for this request
	Slots() []Slot
	// Category returns category obviously
	Category() []Category
	// Type return the publisher type
	Type() ImpressionType
	// Is this publisher accept under floor ads or not ?
	UnderFloor() bool
}
