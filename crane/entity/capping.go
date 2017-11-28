package entity

// Interface interface capping
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
}

// SortByCap is the sort entry based on selected/ad capping/campaign capping/cpm (order is matter)
type SortByCap []Advertise

// Len return the len of array
func (a SortByCap) Len() int {
	return len(a)
}

// Swap two item in array
func (a SortByCap) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less return if the index i is less then index j?
func (a SortByCap) Less(i, j int) bool {
	// This is a multi-sort function.
	iCP := a[i].Capping().View() / a[i].Capping().Frequency()
	jCP := a[j].Capping().View() / a[j].Capping().Frequency()
	if iCP != jCP {
		return iCP < jCP
	}

	//if a[i].Capping().View() != a[j].Capping().View() {
	//	return a[i].Capping().View() < a[j].Capping().View()
	//}

	return a[i].CPM() < a[j].CPM()
}
