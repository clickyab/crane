package entity

// Slot is the slot of the app
type Slot interface {
	// ID of the slot not changed (permanent)
	ID() string

	PublicID() string
	// TrackID of slot
	ReservedHash() string
	// Width return the primary size of this slot
	Width() int
	// Height return the primary size of this slot
	Height() int
	// Size return the clickyab size
	Size() int
	// SetSlotCTR the ctr from database
	SetSlotCTR(float64)
	// SlotCTR the ctr from database
	SlotCTR() float64
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
	// IsSizeAllowed return true if the size is allowed in this slot
	IsSizeAllowed(int, int) bool

	ExtraParams() map[string]string
}
