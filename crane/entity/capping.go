package entity

// Capping interface capping object for all the campaign ads
type Capping interface {
	// View return the view of this campaign for this user
	View() int
	// Frequency return the frequency for this object
	Frequency() int
	// Capping return the frequency capping value, the view/frequency
	// IncView increase the vie
	IncView(int, bool)
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
	// This is a multisort function.
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
