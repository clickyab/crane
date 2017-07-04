package entity

// Platforms consist of app, web, vast
type Platforms string

const (

	// AppPlatform AppPlatform
	AppPlatform Platforms = "app"
	// VastPlatform VastPlatform
	VastPlatform Platforms = "vast"
	// WebPlatform WebPlatform
	WebPlatform Platforms = "web"

	// NativePlatform NativePlatform
	NativePlatform Platforms = "native"
)

// QPublisher will contains all query related to publisher
type QPublisher interface {
	// Find publisher by ID
	Find(int64) (Publisher, error)
	// ByPlatform get Publisher name, Platform and supplier
	ByPlatform(string, Platforms, string) (Publisher, error)
}
