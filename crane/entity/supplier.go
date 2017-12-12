package entity

// Supplier is the ad-network interface
type Supplier interface {
	// Name of Supplier
	Name() string
	// Token of this for web request
	Token() string
	// CPMFloor is the floor for this network. the publisher must be greeter equal to this
	DefaultFloorCPM() int64
	// DefaultSoftFloorCPM is the default, when the site/app is not available
	DefaultSoftFloorCPM() int64
	// DefaultMinBid return the default min bid for this supplier
	DefaultMinBid() int64
	// This publisher bid type
	BidType() BIDType
	// DefaultCTR for this supplier
	DefaultCTR() float64
	// AllowCreate indicated if this supplier can create publisher on demand
	AllowCreate() bool
	// TinyMark means we can add our mark to it
	TinyMark() bool
	// ShowDomain is a domain that all links are generated against it
	ShowDomain() string
}
