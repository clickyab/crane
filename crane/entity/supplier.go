package entity

// Supplier is the ad-network interface
type Supplier interface {
	// ID is the network id in our system
	ID() int64
	// Name of Supplier
	Name() string
	// CPMFloor is the floor for this network. the publisher must be greeter equal to this
	CPMFloor() int64
	// AcceptedTypes is the types that this network can request
	AcceptedTypes() []AdType
	// ExcludedNetwork is the black listed network for this.
	ExcludedDemands() []string
	// CountryWhiteList is the list of country accepted for this
	CountryWhiteList() []Country
}
