package entity

import "time"

// Seat is the slot of the app
type Seat interface {
	// PublicID of the slot not changed (permanent)
	PublicID() string
	// ReservedHash of slot
	ReservedHash() string
	// Width return width of seat
	Width() int
	// Height return height of seat
	Height() int
	// Size return the clickyab size
	Size() int
	// Bid return winner bid
	Bid() float64
	// CPM return the cpm of this seat (after winner and bid is set)
	CPM() float64
	// Set winner ad for this slot, first is bid, last param is cpm
	SetWinnerAdvertise(Advertise, float64, float64)
	// WinnerAdvertise return the winner
	WinnerAdvertise() Advertise
	// ShowURL get the show url usable for async calls
	ShowURL() string
	// SetClickURL is the setter for click url of this ad in slot
	ClickURL() string
	// CTR return current ctr for this size in publisher
	CTR() float64
	// Type of seat
	Type() string
	// SubType of seat
	SubType() string
	// ShowT the fucking iframe injection
	ShowT() bool
	// MinBid is the minimum CPC requested for this requests
	MinBid() int64
	// ImpressionTime is the time of impression (if this is impression seat, on show, its current time,
	// but if this is click, then its the impression time, not current)
	ImpressionTime() time.Time
	// SupplierCPM return cpm from supplier
	SupplierCPM() float64
	// FatFinger return true if we need to prevent sudden click
	FatFinger() bool
	// Some seat need extra filter so the specific ad could be removed, for example in mime-type if a seat
	// only accept video/mp4 but another seat in same imp accept image/jpeg then we can not use normal filters
	Acceptable(Advertise) bool
}

// VastSeat is a seat with vast compatibility
type VastSeat interface {
	Seat
	// Linear is only usable in vast subsystem!
	Linear() bool
	// Duration is the function to handle the vast duration
	Duration() time.Duration
	// SkipAfter duration in vast
	SkipAfter() time.Duration
}
