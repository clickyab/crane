package entity

type (
	// PublisherAttributes is the publisher attributes
	PublisherAttributes string
)

const (
	// PublisherTypeApp is the app
	PublisherTypeApp = 1
	// PublisherTypeWeb is the web
	PublisherTypeWeb = 2
	// PublisherTypeVast is the vast
	PublisherTypeVast = 3
)

// BIDType is the bid type for this imp cpc or cpm
type BIDType string

const (
	// BIDTypeCPC is the cost per click type
	BIDTypeCPC BIDType = "CPC"
	//BIDTypeCPM is the cost per view type
	BIDTypeCPM BIDType = "CPM"
)

// Publisher is the publisher interface
type Publisher interface {
	// GetID return the publisher id
	ID() int64
	// FloorCPM is the floor cpm for publisher
	FloorCPM() int64
	// SoftFloorCPM is the soft version of floor cpm. if the publisher ahs it, then the system
	// try to use this as floor, but if this is not available, the FloorCPM is used
	SoftFloorCPM() int64
	// Name of publisher
	Name() string
	// Active is the publisher active?
	Active() bool
	// Type return the publisher type
	AcceptedTarget() Target
	// Attributes is the generic attribute system
	Attributes(PublisherAttributes) interface{}
	// BIDType return this publisher bid type
	BIDType() BIDType
	// MinCPC is the minimum CPC requested for this requests
	MinCPC() int64
	// AcceptedTypes is the type accepted by this impression
	AcceptedTypes() []AdType
	// Is this publisher accept under floor ads or not ?
	UnderFloor() bool
	// Supplier return the exchange object for this publisher
	Supplier() Supplier
}
