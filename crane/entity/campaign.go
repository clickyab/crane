package entity

// Target is the target of this campaign
type Target int

const (

	// TargetInvalid is the invalid banner in system
	TargetInvalid Target = iota
	// TargetWeb TargetWeb
	TargetWeb
	// TargetApp is the app targeted campaign
	TargetApp
	// TargetVast is the vast target
	TargetVast
	// TargetNative is the native platform
	TargetNative
)

// Campaign is the single campaign in system
type Campaign interface {
	// ID return the campaign id
	ID() int64
	// Name is the campaign name
	Name() string
	// MaxBID get the campaign max bid
	MaxBID() int64
	// Make sure the result is >= 1
	Frequency() int
	// Target return the target of this campaign
	Target() Target
	//BlackListPublisher shows publishers in blacklist
	BlackListPublisher() []string
	//BlackListPublisher shows publishers in blacklist
	WhiteListPublisher() []string
	// AppBrands return campaign app brands
	AppBrands() []string
	// AppCarriers return campaign app carriers
	AppCarriers() []string
	//AllowedOS return os blacklist of a campaign
	AllowedOS() []string
	//Country return country
	Country() []string
	// Province returns province ID
	Province() []string
	//LatLon return LanLon and radius to accept ad
	LatLon() (float64, float64, float64)
	// Category return the category of this campaign
	Category() []Category
	// Attributes return the ad specific attributes
	Attributes() map[string]interface{}
	// WebMobile campaign web mobile on or off
	WebMobile() bool
	// Web campaign web on or off
	Web() bool
	Hoods() []string
	// Isp list of campaign isp(s)
	Isp() []string
	// NetProvider net providers for certain campaign
	NetProvider() []string
}

// IsSizeAllowed return if the size is allowed in target type or not
func (t Target) IsSizeAllowed(w, h int) bool {
	// TODO : Write the entire body
	return true
}
