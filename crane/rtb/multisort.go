package rtb

import (
	"clickyab.com/crane/crane/entity"
)

// byMulti sort by multi sort
type byMulti struct {
	Video bool
	Ads   []entity.Advertise
}

func (a byMulti) Len() int {
	return len(a.Ads)
}
func (a byMulti) Swap(i, j int) {
	a.Ads[i], a.Ads[j] = a.Ads[j], a.Ads[i]
}
func (a byMulti) Less(i, j int) bool {

	if a.Ads[i].Capping().Selected() != a.Ads[j].Capping().Selected() {
		return !a.Ads[i].Capping().Selected()
	}

	if a.Ads[i].Capping().AdCapping(a.Ads[i].ID()) != a.Ads[j].Capping().AdCapping(a.Ads[j].ID()) {
		return a.Ads[i].Capping().AdCapping(a.Ads[i].ID()) < a.Ads[j].Capping().AdCapping(a.Ads[j].ID())
	}

	if a.Video {
		if a.Ads[i].Type() != a.Ads[j].Type() {
			return a.Ads[i].Type() == entity.AdTypeVideo
		}
	}
	return a.Ads[i].CPM() > a.Ads[j].CPM()
}
