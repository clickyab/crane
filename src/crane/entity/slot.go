package entity

// Slot is the slot of the app
type Slot interface {
	// ID of slot
	ID() int64
	// PublicID of slot
	PublicID() string
	// AllowedSize of slot, if the slot allowed multiple size, then the result is mor than one int
	AllowedSize() []int
	// Size return the primary size of this slot
	Size() int
	// StateID is an string for this slot, its a random at first but the value is not changed at all other calls
	StateID() string
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
}
