package entity

import (
	"net/url"
	"time"
)

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
	SetWinnerAdvertise(Creative, float64, float64)
	// WinnerAdvertise return the winner
	WinnerAdvertise() Creative
	// ShowURL get the show url usable for async calls
	ImpressionURL() *url.URL
	// ClickURL is the click url for this advertise
	ClickURL() *url.URL
	// WinRequest is the win request url for this advertise
	WinNoticeRequest() *url.URL
	// CTR return current ctr for this size in publisher
	CTR() float64
	// Type of seat
	Type() InputType
	// RequestType of seat
	RequestType() RequestType
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
	Acceptable(Creative) bool
	// MinCPC return min cpc of seat
	MinCPC() float64
	// The minimum CPM allowed by this, HARD FLOOR, after share calculation
	MinCPM() float64
	// Soft minimum. lower than this means no sec biding. must not be lower than min cpm
	SoftCPM() float64
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

// NativeSeat is the seat for native
type NativeSeat interface {
	Seat
	// Filters return array of slots (only required ones, not all of them)
	Filters() []Filter
}
