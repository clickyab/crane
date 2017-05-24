package entity

// URLProvider is the provider for urls used in the system
type URLProvider interface {
	// ShowURL Create the url for showing the ad as a single ad, its not related to ad
	ShowURL(Slot, Impression, Publisher) string
	// ClickURL return a click link for app
	ClickURL(Slot, Impression, Publisher, Advertise)
}
