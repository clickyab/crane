package capping

import "clickyab.com/crane/crane/entity"

// sortByCap is the sort entry based on selected/ad capping/campaign capping/cpm (order is matter)
type sortByCap []entity.Advertise

// Len return the len of array
func (a sortByCap) Len() int {
	return len(a)
}

// Swap two item in array
func (a sortByCap) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less return if the index i is less then index j?
func (a sortByCap) Less(i, j int) bool {
	// This is a multi-sort function.
	iCP := a[i].Capping().View() / a[i].Capping().Frequency()
	jCP := a[j].Capping().View() / a[j].Capping().Frequency()
	if iCP != jCP {
		return iCP < jCP
	}

	return a[i].Capping().View() < a[j].Capping().View()
}
