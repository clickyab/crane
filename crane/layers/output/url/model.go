package url

import (
	"time"

	"clickyab.com/crane/crane/entity"
)

// Data is base interface for processing click event
type Data struct {
	// Target url for ad
	Target string
	// ImpressionTrackID is for tracking the impression
	ImpressionTrackID string
	// Publisher is the publisher name
	Publisher string
	// ClientID is unique user id
	ClientID string
	// Supplier is the name of supplier
	Supplier string
	// IP of client
	IP string
	// UserAgent of client
	UserAgent string
	// SlotTrackID is unique track id for slot
	SlotTrackID string
	// WinnerBID is cost of ad
	WinnerBID int64
	// AdvertiseID is the ad id
	AdvertiseID int64
	// GenTime is when url has been generated
	GenTime time.Time
	// ClickTime is when click happened
	ClickTime time.Time
	// Status of click which can be one of clickStatus
	Status entity.ClickStatus
}
