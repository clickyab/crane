package entity

import (
	"fmt"
)

// SortStrategyType is kind of sort strategy <ctr,cpm>
// @Enum{
// }
type SortStrategyType string

const (
	// SortByCTR set sort strategy to ctr
	SortByCTR SortStrategyType = "ctr"
	// SortByCPM set sort strategy to cpm
	SortByCPM SortStrategyType = "cpm"
)

// SortableCreative sort by multi sort
type SortableCreative struct {
	Video        bool
	SortStrategy SortStrategyType
	Ads          []SelectedCreative
}

func (a SortableCreative) Len() int {
	return len(a.Ads)
}
func (a SortableCreative) Swap(i, j int) {
	a.Ads[i], a.Ads[j] = a.Ads[j], a.Ads[i]
}
func (a SortableCreative) Less(i, j int) bool {

	if a.Ads[i].Capping().Selected() != a.Ads[j].Capping().Selected() {
		return !a.Ads[i].Capping().Selected()
	}

	if a.Ads[i].Capping().AdCapping(a.Ads[i].ID()) != a.Ads[j].Capping().AdCapping(a.Ads[j].ID()) {
		return a.Ads[i].Capping().AdCapping(a.Ads[i].ID()) < a.Ads[j].Capping().AdCapping(a.Ads[j].ID())
	}

	if a.Video {
		if a.Ads[i].Type() != a.Ads[j].Type() {
			return a.Ads[i].Type() == AdTypeVideo
		}
	}

	if a.SortStrategy == SortByCTR {
		if a.Ads[i].CalculatedCTR() != a.Ads[j].CalculatedCTR() {
			return a.Ads[i].CalculatedCTR() > a.Ads[j].CalculatedCTR()
		}
	}

	if a.Ads[i].CalculatedCPM() != a.Ads[j].CalculatedCPM() {
		return a.Ads[i].CalculatedCPM() > a.Ads[j].CalculatedCPM()
	}

	return a.Ads[i].CalculatedCPC() > a.Ads[j].CalculatedCPC()
}

// String is a helper method for debugging
func (a SortableCreative) String() string {
	res := "---\n"
	for i := range a.Ads {
		res += fmt.Sprintf("CP:%d AD:%d CPM:%v CPC:%v ADC:%v SEL:%v\n",
			a.Ads[i].Campaign().ID(),
			a.Ads[i].ID(),
			a.Ads[i].CalculatedCPM(),
			a.Ads[i].CalculatedCPC(),
			a.Ads[i].Capping().AdCapping(a.Ads[i].ID()),
			a.Ads[i].Capping().Selected(),
		)
	}
	res += "===\n"
	return res
}