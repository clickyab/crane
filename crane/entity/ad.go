package entity

import "net/url"

// AdType is the type supported by ads
type AdType string

const (
	// AdTypeBanner is the banner type
	AdTypeBanner AdType = "banner"
	// AdTypeDynamic is the dynamic type. the code is html
	AdTypeDynamic AdType = "dyn"
	// AdTypeVideo is the video type
	AdTypeVideo AdType = "video"
	// AdTypeHTML is the html ad type
	AdTypeHTML AdType = "html"
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
	SetWinnerBID(int64)
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
	// Media advertise such as image,gif,video,...
	Media() string
	// TargetURL end of show ad redirect to this URL
	TargetURL() url.URL
	// Attributes return the ad specific attributes
	Attributes() map[string]interface{}
	// Duplicate is a hackish function to handle the duplicate of interface
	Duplicate() Advertise
}
