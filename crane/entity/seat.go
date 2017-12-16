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
	// CPM return the cpm of this seat (after winner and bid is set)
	CPM() float64
	// Set winner ad for this slot, first is bid, last param is cpm
	SetWinnerAdvertise(Advertise, float64, float64)
	// WinnerAdvertise return the winner
	WinnerAdvertise() Advertise
	// ShowURL get the show url usable for async calls
	ShowURL() string
	// SetClickURL is the setter for click url of this ad in slot
	ClickURL() string
	// CTR return current ctr for this size in publisher
	CTR() float64
	// Type of seat
	Type() string
	// ShowT the fucking iframe injection
	ShowT() bool
	// MinBid is the minimum CPC requested for this requests
	MinBid() int64
}
