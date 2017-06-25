package entity

// ClickProvider is the provider for urls used in the system
type ClickProvider interface {
	// ClickURL return a click link for app
	ClickURL(Slot, Impression, Publisher, Advertise)
}
