package entity

// RequestType is the request type
type RequestType string

const (
	// RequestTypeWeb is web type (for clickyab web)
	RequestTypeWeb RequestType = "web"
	// RequestTypeApp is app type (for clickyab sdk only)
	RequestTypeApp RequestType = "app"
	// RequestTypeNative is native type (for clickyab native only)
	RequestTypeNative RequestType = "native"
	// RequestTypeVast is vast type (for clickyab vast only)
	RequestTypeVast RequestType = "vast"
	// RequestTypeDemand is demand type
	RequestTypeDemand RequestType = "demand"
)

// Context is the single impression object
type Context interface {
	// Type return the request type.
	Type() RequestType
	// Request data comes from request for every user
	// like ip,user agent,client id,...
	Request
	// Publisher return the publisher that this client is going into system from that
	Publisher() Publisher
	// Slots is the slot for this request
	Seats() []Seat
	// Category returns category obviously
	Category() []Category
	// User return user data
	User() User
	// Tiny means that the logo of clickyab should be shown (true) or not
	Tiny() bool
	// Currency
	Currency() string
	// MultiVideo determine this request can have multiple video
	MultiVideo() bool
	// FloorDiv is floor-cpm divider
	FloorDiv() int64
	// Capping is required otr not
	Capping() bool
}
