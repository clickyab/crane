package exchange

// Advertise is the single advertise interface
type Advertise interface {
	// GetID return the id of advertise
	ID() string
	// MaxCPM return the max cpm of this ad, from the deman
	MaxCPM() int64
	// Width return the size
	Width() int
	// Height return the size
	Height() int
	// return the url to call for show
	URL() string
	Rater

	AdvertiseExtra
}

// AdvertiseExtra is the parameters from our system passed to advertise
type AdvertiseExtra interface {
	// Return the track id, it must be randomly generated code and after the first call
	// must not change in a one call
	TrackID() string
	// SetWinnerCPM is the
	SetWinnerCPM(int64)
	// WinnerCPM return the winner value already set on SetWinnerCPM zero if not set already
	WinnerCPM() int64
	// Demand return the demand registered with this ad
	Demand() Demand
}
