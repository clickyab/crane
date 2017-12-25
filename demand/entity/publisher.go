package entity

// PublisherType is the type of publisher
type PublisherType int

const (
	// PublisherTypeApp is the app
	PublisherTypeApp PublisherType = 1
	// PublisherTypeWeb is the web
	PublisherTypeWeb PublisherType = 2
)

// BIDType is the bid type for this imp cpc or cpm
type BIDType string

const (
	// BIDTypeCPC is the cost per click type
	BIDTypeCPC BIDType = "CPC"
	//BIDTypeCPM is the cost per view type
	BIDTypeCPM BIDType = "CPM"
)

// PublisherAttributes is the key for publisher attributes
type PublisherAttributes int

const (
	// PAMobileAd determine if banner ad show be shown in mobile version or not
	PAMobileAd PublisherAttributes = iota
	// PAFatFinger determine sensitivity of touch ads in app
	PAFatFinger
)

// Publisher is the publisher interface
type Publisher interface {
	ID() int64
	// FloorCPM is the floor cpm for publisher
	FloorCPM() int64
	// SoftFloorCPM is the soft version of floor cpm. if the publisher ahs it, then the system
	// try to use this as floor, but if this is not available, the FloorCPM is used
	SoftFloorCPM() int64
	// Name of publisher
	Name() string
	// BIDType return this publisher bid type
	BIDType() BIDType
	// MinBid is the minimum CPC requested for this requests
	MinBid() int64
	// Supplier return the exchange object for this publisher
	Supplier() Supplier
	// CTR returns ctr of a slot with specific size
	CTR(int) float64
	// Type return type of this publisher
	Type() PublisherType
	// Attributes si any other attributes that is not generally required for other part of the system
	Attributes() map[PublisherAttributes]interface{}
}