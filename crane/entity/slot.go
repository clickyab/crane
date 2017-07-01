package entity

// Slot is the slot of the app
type Slot interface {
	ID() string
	// TrackID of slot
	TrackID() string
	// Width return the primary size of this slot
	Width() int
	// Height return the primary size of this slot
	Height() int
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
	// SetClickURL set internal system url for our logic and after redirect to TargetURL Advertise
	SetClickURL() string
	// ClickURL is internal system url for our logic and after redirect to TargetURL Advertise
	ClickURL() string
	// IsSizeAllowed return true if the size is allowed in this slot
	IsSizeAllowed(int, int) bool
	// Attribute return meta data for slot
	Attribute() map[string]interface{}
}
