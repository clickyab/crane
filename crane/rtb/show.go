package rtb

import (
	"sort"

	"clickyab.com/crane/crane/capping"
	"clickyab.com/crane/crane/entity"
)

// select sets winner ad detail into context
func Select(c entity.Context, ads map[int][]entity.Advertise) {
	ads = capping.GetCapping(c.User().ID(), ads, "", c.Slots()...)
	for i := range c.Slots() {
		slot := c.Slots()[i]
		winnerAd, winnerBid := bidding(c, ads[slot.Size()])

		winnerAd.SetWinnerBID(winnerBid, true)
		slot.SetWinnerAdvertise(winnerAd)

		capping.StoreCapping(c.User().ID(), winnerAd.ID())
	}
}

// TODO: care for no ad
// multipleVideo is false hardcoded
// return winner bid with its bid
func bidding(c entity.Context, ads []entity.Advertise) (ad entity.Advertise, wBid int64) {

	var (
		publisher  = c.Publisher()
		exeedFloor []entity.Advertise
		underFloor []entity.Advertise

		sortingData byMulti
	)

	for i := range ads {
		if publisher.FloorCPM() > ads[i].CPM() {
			exeedFloor = append(exeedFloor, ads[i])
		} else {
			underFloor = append(underFloor, ads[i])
		}
	}

	if len(exeedFloor) > 0 {
		sortingData = byMulti{Video: false, Ads: exeedFloor}
	} else {
		sortingData = byMulti{Video: false, Ads: underFloor}
	}

	sort.Sort(sortingData)
	sorted := sortingData.Ads
	ad = sorted[0]

	if len(exeedFloor) > 0 {
		cpm := getSecondCPM(publisher.FloorCPM(), exeedFloor)
		wBid = winnerBid(cpm, ad.CTR())
	} else {
		wBid = winnerBid(int64(ad.CPM()), ad.CTR())
	}

	return
}
