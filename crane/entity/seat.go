package entity

// Seat is the slot of the app
type Seat interface {
	// PublicID of the slot not changed (permanent)
	PublicID() string
	// ReservedHash of slot
	ReservedHash() string
	// Size return the clickyab size
	Size() int
	// SetBid set bid for winner
	SetBid(float64)
	// Bid return winner bid
	Bid() float64
	// Set winner ad for this slot
	SetWinnerAdvertise(Advertise)
	// WinnerAdvertise return the winner
	WinnerAdvertise() Advertise
	// SetShowURL set the show url usable for async calls
	SetShowURL(string)
	// ShowURL get the show url usable for async calls
	ShowURL() string
	// SetClickURL is the setter for click url of this ad in slot
	SetClickURL(string)
	// ClickURL is the setter for click url of this ad in slot
	ClickURL() string
}
