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
	// AdTypeNative is the native type
	AdTypeNative AdType = 4
)

// Creative is the single advertise interface
type Creative interface {
	// GetID return the id of advertise
	ID() int64
	// Type return the type of ad
	Type() AdType
	// Campaign return the ad campaign
	Campaign() Campaign
	// AdCTR the ad ctr from database (its not calculated from )
	AdCTR() float64
	//// MaxBID get the campaign max bid
	MaxBID() int64
	//// Target return the target of this campaign
	Target() Target
	// Size returns ads size
	Size() int
	// Width return the width
	Width() int
	// Height return the height of banner
	Height() int
	// Duration of the ad if it have meaning. normally usable for vast, in second
	// TODO : duration could be removed from here and moved to attributes, if there is no other
	// need to duration, then its safe to move it to attributes
	Duration() int
	// Capping return the current capping object
	Capping() Capping
	// SetCapping set the current capping
	SetCapping(Capping)
	// Attributes return the ad specific attributes
	Attributes() map[string]interface{}
	// Media return image of ad
	// Deprecated: use Asset function
	Media() string
	// Target return ad target of the target
	TargetURL() string
	// TODO: remove this later
	// CampaignAdID return campaign_ad primary
	CampaignAdID() int64
	// MimeType of media
	MimeType() string
	// Asset return the asset that pass all filters and type is exactly matched the value
	Assets(AssetType, int, ...AssetFilter) []Asset
}

// SelectedCreative is used during the rtb loop.
// sometimes I hate the Go type system
type SelectedCreative interface {
	Creative
	CalculatedCTR() float64
	CalculatedCPM() int64
	CalculatedCPC() int64
	IsSecBid() bool
}
