package entity

// Seat is the slot of the app
type Seat interface {
	// PublicID of the slot not changed (permanent)
	PublicID() string
	// ReservedHash of slot
	ReservedHash() string
	// Width return width of seat
	Width() int
	// Height return height of seat
	Height() int
	// Size return the clickyab size
	Size() int
	// Bid return winner bid
	Bid() float64
	// Set winner ad for this slot
	SetWinnerAdvertise(Advertise, float64)
	// WinnerAdvertise return the winner
	WinnerAdvertise() Advertise
	// ShowURL get the show url usable for async calls
	ShowURL() string
	// SetClickURL is the setter for click url of this ad in slot
	ClickURL() string
	// Supplier return supplier
	Supplier() string
	// CTR return current ctr for this size in publisher
	CTR() float64
}
