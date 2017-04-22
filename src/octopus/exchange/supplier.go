package exchange

// Supplier is the ad-network interface
type Supplier interface {
	// Name of Supplier
	Name() string
	// CPMFloor is the floor for this network. the publisher must be greeter equal to this
	FloorCPM() int64
	// SoftFloorCPM is the soft version of floor cpm. if the publisher ahs it, then the system
	// try to use this as floor, but if this is not available, the FloorCPM is used
	SoftFloorCPM() int64
	// ExcludedNetwork is the black listed network for this.
	ExcludedDemands() []string
	// CountryWhiteList is the list of country accepted for this
	CountryWhiteList() []Country
	// Share return the share of this supplier
	Share() int
	// Renderer return the renderer of this supplier
	Renderer() Renderer
}
