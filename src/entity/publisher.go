package entity

type (

	// PublisherAttributes is the publisher attributes
	PublisherAttributes string
)

// Publisher is the publisher interface
type Publisher interface {
	// Name of publisher
	Name() string
	// FloorCPM is the floor cpm for publisher
	FloorCPM() int64
	// SoftFloorCPM is the soft version of floor cpm. if the publisher ahs it, then the system
	// try to use this as floor, but if this is not available, the FloorCPM is used
	SoftFloorCPM() int64
	// Active is the publisher active?
	Active() bool
	// Attributes is the generic attribute system
	Attributes(PublisherAttributes) interface{}
	// AcceptedTypes is the type accepted by this impression
	AcceptedTypes() []AdType
	// Is this publisher accept under floor ads or not ?
	UnderFloor() bool
	// Supplier return the exchange object for this publisher
	Supplier() Supplier
}
