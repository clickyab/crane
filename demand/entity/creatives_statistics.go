package entity

// CreativeStatistics is the single advertise interface
type CreativeStatistics interface {
	// CreativeType return the type of ad
	CreativeType() AdType
	// TotalCount return total count of creatives base on type
	TotalCount() int64
	// MinCTR return network creatives statistics
	MinCTR() float64
	// MaxCTR return network creatives statistics
	MaxCTR() float64
	// AvgCTR return network creatives statistics
	AvgCTR() float64
}
