package entity

// Capping interface capping
type Capping interface {
	// View return the view of this campaign for this user
	View() int
	// View return the view of this campaign for this user
	AdView(int64) int
	// Frequency return the frequency for this user
	Frequency() int
	// Capping return the frequency capping value, the view/frequency
	Capping() int
	// Capping return the frequency capping value, the view/frequency
	AdCapping(int64) int
	// IncView increase the vie
	IncView(int64, int, bool)
	// Selected return if this campaign is already selected in this batch
	Selected() bool
	// Store the capping
	Store(int64)
}
