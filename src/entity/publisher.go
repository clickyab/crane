package entity

type (
	// PublisherType is the publisher type
	PublisherType string
	// PublisherAttributes is the publisher attributes
	PublisherAttributes string
)

const (
	// PublisherTypeApp is the app
	PublisherTypeApp PublisherType = "app"
	// PublisherTypeWeb is the web
	PublisherTypeWeb PublisherType = "web"
	// PublisherTypeVast is the vast
	PublisherTypeVast PublisherType = "vast"
)

// Publisher is the publisher interface
type Publisher interface {
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
	Type() PublisherType
	// Attributes is the generic attribute system
	Attributes(PublisherAttributes) interface{}
	// MinCPC is the minimum CPC requested for this requests
	MinCPC() int64
	// AcceptedTypes is the type accepted by this impression
	AcceptedTypes() []AdType
	// Is this publisher accept under floor ads or not ?
	UnderFloor() bool
	// Supplier return the exchange object for this publisher
	Supplier() Supplier
}

