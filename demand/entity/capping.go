package entity

// Capping interface capping
type Capping interface {
	// View return the view of this campaign for this user
	View() int32
	// View return the view of this campaign for this user
	AdView(int32) int32
	// Frequency return the frequency for this user
	Frequency() int32
	// Capping return the frequency capping value, the view/frequency
	Capping() int32
	// Capping return the frequency capping value, the view/frequency
	AdCapping(int32) int32
	// IncView increase the vie
	IncView(int32, int32, bool)
	// Selected return if this campaign is already selected in this batch
	Selected() bool
}
