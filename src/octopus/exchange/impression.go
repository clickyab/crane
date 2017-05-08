package exchange

import (
	"net"
	"time"
)

// ImpressionPlatform is the publisher type
type ImpressionPlatform int

var impressionPlatformName = [...]string{"app", "vast", "web"}

// String return string name of platform
func (m ImpressionPlatform) String() string {
	return impressionPlatformName[m]
}

const (
	// ImpressionPlatformApp is the app
	ImpressionPlatformApp ImpressionPlatform = iota
	// ImpressionPlatformVast is the vast
	ImpressionPlatformVast
	// ImpressionPlatformWeb is the web
	ImpressionPlatformWeb
)

// ImpressionAttributes is the imp attr key
type ImpressionAttributes string

// Impression is the single impression object
type Impression interface {
	// TrackID return the random id of this imp object
	TrackID() string
	// IP return the client ip
	IP() net.IP
	// UserAgent return the client user agent
	UserAgent() string
	// Source return the publisher that this client is going into system from that
	Source() Publisher
	// Location of the request
	Location() Location
	// Attributes is the generic attribute system
	Attributes() map[string]interface{}
	// Slots is the slot for this request
	Slots() []Slot
	// Category returns category obviously
	Category() []Category
	// Platform return the publisher type
	Platform() ImpressionPlatform
	// Is this publisher accept under floor ads or not ?
	UnderFloor() bool
	// Time time of the impression
	Time() time.Time
}
