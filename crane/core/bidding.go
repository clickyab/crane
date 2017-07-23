package core

import (
	"clickyab.com/crane/crane/entity"
	"github.com/Sirupsen/logrus"
)

// Select is the real time bidding on ads
func Select(ads []entity.Advertise, impression entity.Impression, slotIndex int) {
	floorCPM := impression.Publisher().FloorCPM()
	slot := impression.Slots()[slotIndex]

	if len(ads) < 2 {
		if len(ads) == 0 {
			logrus.Warn("no ad passed by for bidding")
			return
		}

		slot.SetWinnerAdvertise(ads[0])
		if impression.Publisher().FloorCPM() > ads[0].CPM() {
			ads[0].SetWinnerBID(ads[0].CPM())
			return
		}

		ads[0].SetWinnerBID(floorCPM + 1)
		return
	}

	capping := GetCapping(impression.ClientID(), ads)
	ads = capping.Sort()

	firstAd := ads[0]
	secondAd := ads[1]

	slot.SetWinnerAdvertise(firstAd)
	capping.IncView(firstAd)
	if firstAd.CPM() < floorCPM {
		firstAd.SetWinnerBID(firstAd.CPM())
		return
	} else if secondAd.CPM() < floorCPM {
		firstAd.SetWinnerBID(floorCPM + 1)
		return
	}

	firstAd.SetWinnerBID(secondAd.CPM() + 1)
}
