package entity

// AdType is the type supported by ads
type AdType string

const (
	// AdTypeBanner is the banner type
	AdTypeBanner AdType = "banner"
	// AdTypeVideo is the video type
	AdTypeVideo AdType = "video"
	// AdTypeHTML is the html ad type
	AdTypeHTML AdType = "html"
)

// Advertise is the single advertise interface
type Advertise interface {
	// GetID return the id of advertise
	ID() string
	// Type return the type of ad
	Type() AdType
	// SetCPM set the cpm for this ad in the system after select
	CPM() int64
	// Width return the size
	Width() int
	// Height return the size
	Height() int
	// return the code to show
	Code() string
	// Win is called when the ad is selected for a slot. it must call the provider
	Win()
}
