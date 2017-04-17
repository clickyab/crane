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
	// Attributes is the generic attribute system
	Attributes() map[string]interface{}
	// Supplier return the exchange object for this publisher
	Supplier() Supplier
}
