package entity

// InputType is demand
type InputType string

// InputTypeDemand is demand
const InputTypeDemand InputType = "demand"

func (r InputType) String() string {
	return string(r)
}

// RequestType is the request type
type RequestType string

// IsValid will return true if request if true
func (r RequestType) IsValid() bool {
	return r == RequestTypeBanner || r == RequestTypeVast || r == RequestTypeNative
}

const (
	// RequestTypeNative is native type (for clickyab native only)
	RequestTypeNative RequestType = "native"
	// RequestTypeVast is vast type (for clickyab vast only)
	RequestTypeVast RequestType = "vast"
	// RequestTypeBanner is vast type (for clickyab vast only)
	RequestTypeBanner RequestType = "banner"
)

func (r RequestType) String() string {
	return string(r)
}

// CappingMode decide how to handle capping
type CappingMode int

// Validate capping modes
func (c CappingMode) Validate() bool {
	return c == CappingStrict || c == CappingReset || c == CappingNone
}

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
	Type() InputType
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
	// Rate return ratio currency conversion to IRR
	Rate() float64
	// ConnectionType return network for app 2g,3g,4g,...
	ConnectionType() int
	// Carrier return carrier for app
	Carrier() string
	// Brand return brand for app
	Brand() string
	// FatFinger return true if we need to prevent sudden click
	FatFinger() bool
	// PreventDefault is a way to handle old sdk version
	PreventDefault() bool
	// Strategy of biding (cpm, cpc)
	Strategy() Strategy
	// UnderFloor means that this supplier allow to pass underfloor value.
	// normally used only for clickyab
	UnderFloor() bool
	// TV true if true view
	TV() bool
	// BannerMarkup true if do not need iframe
	BannerMarkup() bool
	// GetCreativesStatistics return statistics of all active network creatives base on it's type
	GetCreativesStatistics() []CreativeStatistics
}
