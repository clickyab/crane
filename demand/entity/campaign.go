package entity

import "strings"

// Target is the target of this campaign
type Target int

const (
	// TargetWeb TargetWeb
	TargetWeb Target = iota
	// TargetApp is the app targeted campaign
	TargetApp
	// TargetVast is the vast target
	TargetVast
	// TargetNative is the native platform
	TargetNative
)

// Strategy of campaign
type Strategy uint64

var stringStrategy = map[string]Strategy{
	"cpm": StrategyCPM,
	"cpc": StrategyCPC,
}

const (
	// StrategyCPM for campaign
	StrategyCPM Strategy = 1 << iota
	// StrategyCPC for campaign
	StrategyCPC
)

// GetStrategy return strategy from string name (cpc, cpm)
func GetStrategy(s []string) Strategy {
	var res Strategy
	for _, value := range s {
		if v, ok := stringStrategy[strings.ToLower(value)]; ok {
			res |= v
		}
	}
	return res
}

// IsSubsetOf return true if target contain all strategy
func (s Strategy) IsSubsetOf(t Strategy) bool {
	return t|s == t
}

// Valid return true if strategy is valid
func (s Strategy) Valid() bool {
	return (StrategyCPM|StrategyCPC)|s == (StrategyCPM | StrategyCPC)
}

// Campaign is the single campaign in system
type Campaign interface {
	// ID return the campaign id
	ID() int64
	// Name is the campaign name
	Name() string
	// Make sure the result is >= 1
	Frequency() int
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
	LatLon() (bool, float64, float64, float64)
	// Category return the category of this campaign
	Category() []Category
	// WebMobile campaign web mobile on or off
	// @deprecated in favor of attributes TODO : remove this
	WebMobile() bool
	// Web campaign web on or off
	// @deprecated in favor of attributes TODO : remove this
	Web() bool
	// Hoods neighborhood
	// @deprecated do not use!!
	Hoods() []string
	// ISP list of campaign isp(s)
	ISP() []string
	// NetProvider net providers for certain campaign
	NetProvider() []string
	// Strategy can be cpm, cpc
	Strategy() Strategy
}

// IsSizeAllowed return if the size is allowed in target type or not
func (t Target) IsSizeAllowed(w, h int) bool {
	// TODO : Write the entire body
	return true
}
