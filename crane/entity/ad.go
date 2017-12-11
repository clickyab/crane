package entity

type (
	// AdType is the type supported by ads
	AdType int
)

const (
	// AdTypeBanner is the banner type
	AdTypeBanner AdType = 0
	// AdTypeDynamic is the dynamic type. the code is html
	AdTypeDynamic AdType = 2
	// AdTypeVideo is the video type
	AdTypeVideo AdType = 3
	// AdTypeHTML is the html ad type
	AdTypeHTML AdType = 1
	// AdTypeNative is the native ad type
	AdTypeNative AdType = 4
)

// Advertise is the single advertise interface
type Advertise interface {
	// GetID return the id of advertise
	ID() int64
	// Type return the type of ad
	Type() AdType
	// Campaign return the ad campaign
	Campaign() Campaign
	// AdCTR the ad ctr from database (its not calculated from )
	AdCTR() float64
	// Size returns ads size
	Size() int
	// Width return the width
	Width() int
	// Height return the height of banner
	Height() int
	// Capping return the current capping object
	Capping() Capping
	// SetCapping set the current capping
	SetCapping(Capping)
	// Attributes return the ad specific attributes
	Attributes() map[string]interface{}
	// Media return image of ad
	Media() string
	// Target return ad target of the target
	Target() string
	// TODO: remove this later
	// CampaignAdID return campaign_ad primary
	CampaignAdID() int64
}
