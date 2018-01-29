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

func (r RequestType) String() string {
	return string(r)
}

// CappingMode decide how to handle capping
type CappingMode int

const (
	// CappingStrict means we have capping and the capping is strictly used
	CappingStrict CappingMode = iota
	// CappingReset means we have capping but we reset it when the capping is full
	CappingReset
	// CappingNone means no capping at all
	CappingNone
)

// Context is the single impression object
type Context interface {
	// Type return the request type.
	Type() RequestType
	// SubType is the request sub type. normally used for demand
	SubType() RequestType
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
	FloorPercentage() int64
	// Capping is required otr not
	Capping() CappingMode
	// MinBIDPercentage is a hack to handle min bid on multiple types.
	// for example for native its 50%, in new design make sure every
	// type has its own minbid and drop this hack
	MinBIDPercentage() int64
	// SoftFloorCPM is the soft version of floor cpm. if the publisher ahs it, then the system
	// try to use this as floor, but if this is not available, the FloorCPM is used
	SoftFloorCPM() int64
	// BIDType return this publisher bid type
	BIDType() BIDType
	// Rate return ratio currency conversion to IRR
	Rate() float64
	// Network return network for app
	Network() string
	// Carrier return carrier for app
	Carrier() string
	// Brand return brand for app
	Brand() string
	// FatFinger return true if we need to prevent sudden click
	FatFinger() bool
	// PreventDefault is a way to handle old sdk version
	PreventDefault() bool
}
