package entity

type (
	// AdType is the type supported by ads
	AdType string
)

const (
	// AdTypeBanner is the banner type
	AdTypeBanner AdType = "banner"
	// AdTypeDynamic is the dynamic type. the code is html
	AdTypeDynamic AdType = "dyn"
	// AdTypeVideo is the video type
	AdTypeVideo AdType = "video"
	// AdTypeHTML is the html ad type
	AdTypeHTML AdType = "html"
	// AdTypeNative is the native ad type
	AdTypeNative AdType = "native"
)

// Advertise is the single advertise interface
type Advertise interface {
	// GetID return the id of advertise
	ID() int64
	// Type return the type of ad
	Type() AdType
	// Campaign return the ad campaign
	Campaign() Campaign
	// SetCPM set the cpm for this ad in the system after select
	SetCPM(int64)
	// CPM return the current cpm
	CPM() int64
	// SetWinnerBID set the winner bid for this ad if the add is the winner
	// the bool parameter means that the capping must increase
	SetWinnerBID(int64, bool)
	// WinnerBID return the winner bid
	WinnerBID() int64
	// AdCTR the ad ctr from database (its not calculated from )
	AdCTR() float64
	// SetCTR set the calculated CTR
	SetCTR(float64)
	// CTR get the calculated CTR
	CTR() float64
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
	// Duplicate is a hackish function to handle the duplicate of interface
	Duplicate() Advertise
	// Media asd
	Media() string
	// TargetURL asd
	TargetURL() string

	SetSlot(Slot)

	Slot() Slot
}
