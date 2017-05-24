package rtb

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"clickyab.com/exchange/crane/entity"
	"clickyab.com/exchange/services/assert"
	"clickyab.com/exchange/services/eav"
	"clickyab.com/exchange/services/store"
)

const (
	// Mega is the mega store prefix
	Mega string = "MEGA_"
	// MegaIP is the ip subkey
	MegaIP string = "IP"
	// MegaUserAgent is the user agent subkey
	MegaUserAgent string = "UA"
	// MegaPubID is the publisher id subkey
	MegaPubID string = "PID"
	// MegaTimeUnix is the impression timestamp subkey
	MegaTimeUnix string = "TU"
	// MegaAdvertise is the selected ad subkey
	MegaAdvertise string = "AD"
	// MegaSlot is the slot subkey
	MegaSlot string = "SLOT"
)

type resultStore struct {
	// the ads that pass the soft floor
	soft []entity.Advertise
	// the ads that pass the hard floor
	hard []entity.Advertise
	// the ads that pass no floor at all
	under []entity.Advertise
}

func (s resultStore) getData(pub entity.Publisher) (ef entity.SortByCap, floor int64, secBid bool) {
	secBid = true
	floor = pub.SoftFloorCPM()
	// order is to get data from exceed flor, then capping passed and if the config allowed,
	// use the under floor. for under floor there is no second biding pricing
	if len(s.soft) > 0 {
		ef = entity.SortByCap(s.soft)
	} else if len(s.hard) > 0 {
		ef = entity.SortByCap(s.hard)
		floor = pub.FloorCPM()
	} else if pub.UnderFloor() && len(s.under) > 0 {
		ef = entity.SortByCap(s.under)
		secBid = false
	}

	return
}

func createMegaStore(imp entity.Impression) eav.Kiwi {
	kiwi := eav.NewEavStore(Mega + imp.MegaIMP())
	assert.Nil(kiwi.SetSubKey(MegaIP, imp.IP().String()).
		SetSubKey(MegaUserAgent, imp.UserAgent()).
		SetSubKey(MegaPubID, fmt.Sprint(imp.Source().ID())).
		SetSubKey(MegaTimeUnix, fmt.Sprint(time.Now().Unix())).
		Save(*megaImpressionTTL))
	return kiwi
}

// selectCTR is the key function to select an ad for an imp base on real time biding
func selectCTR(
	pCtx context.Context,
	store store.Interface,
	imp entity.Impression,
	ads map[int][]entity.Advertise,
	ch chan map[string]entity.Advertise,
) {

	ctx, cnl := context.WithCancel(pCtx)
	// call cancel on exiting this function so the done channel is fired
	defer cnl()
	// TODO : better implementation
	multiVideo := imp.Source().AcceptedTarget() == entity.TargetVast

	// Get the capping
	slots := imp.Slots()
	pub := imp.Source()
	ads = getCapping(ctx, imp.ClientID(), ads, slots)
	kiwi := createMegaStore(imp)
	wg := sync.WaitGroup{}
	for i := range slots {
		var (
			size    = slots[i].Size()
			noVideo bool
		)

		s := adLoop(ads[size], pub, slots[i], noVideo)

		var sorted []entity.Advertise
		ef, floor, secBid := s.getData(pub)

		if len(ef) == 0 {
			// TODO : Warnings
			store.Push(slots[i].StateID(), "", time.Hour)
			continue
		}

		sort.Sort(ef)
		sorted = []entity.Advertise(ef)
		// Do not do second biding pricing on this ads, they can not pass CPMFloor
		if secBid {
			secondCPM := getSecondCPM(floor, sorted)
			sorted[0].SetWinnerBID(winnerBid(secondCPM, sorted[0].CTR()))
		} else {
			sorted[0].SetWinnerBID(sorted[0].Campaign().MaxBID())
		}

		// Force price on min CPC
		if sorted[0].WinnerBID() < imp.Source().MinCPC() {
			sorted[0].SetWinnerBID(imp.Source().MinCPC())
		}

		sorted[0].Capping().IncView(1, true)
		slots[i].SetWinnerAdvertise(sorted[0])

		if !multiVideo {
			noVideo = noVideo || sorted[0].Type() == entity.AdTypeVideo
		}

		kiwi.SetSubKey(fmt.Sprintf("%s_%d", MegaAdvertise, sorted[0].ID()), fmt.Sprint(sorted[0].WinnerBID()))
		kiwi.SetSubKey(fmt.Sprintf("%s_%d", MegaSlot, sorted[0].ID()), fmt.Sprint(slots[i].ID()))
		assert.Nil(kiwi.Save(*megaImpressionTTL))

		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- map[string]entity.Advertise{slots[i].StateID(): sorted[0]}
			store.Push(slots[i].StateID(), fmt.Sprint(sorted[0].ID()), time.Hour)
		}()
	}

	wg.Wait()
	close(ch)
}

func adLoop(ads []entity.Advertise, pub entity.Publisher, slot entity.Slot, noVideo bool) resultStore {
	res := resultStore{}
	for _, ad := range ads {
		if ad.Type() == entity.AdTypeVideo && noVideo {
			continue
		}
		hard, soft := doBid(ad, pub, slot)
		if ad.WinnerBID() == 0 && soft {
			res.soft = append(res.soft, ad)
		} else if ad.WinnerBID() == 0 && hard {
			res.hard = append(res.hard, ad)
		} else if ad.WinnerBID() == 0 {
			res.under = append(res.under, ad)
		}
	}

	return res
}

func doBid(ad entity.Advertise, pub entity.Publisher, slot entity.Slot) (bool, bool) {
	ad.SetCTR(calculateCTR(
		ad,
		slot,
	))
	ad.SetCPM(cpm(ad.Campaign().MaxBID(), ad.CTR()))
	//exceed cpm floor
	return ad.CPM() >= pub.FloorCPM(), ad.CPM() >= pub.SoftFloorCPM()
}

// CalculateCtr calculate ctr
func calculateCTR(ad entity.Advertise, slot entity.Slot) float64 {
	return (ad.AdCTR()*float64(*adCtrEffect)/100 + slot.SlotCTR()*float64(*slotCtrEffect)/100) / float64(100)
}

//Cpm calculate cpm
func cpm(bid int64, ctr float64) int64 {
	return int64(float64(bid) * ctr * 10.0)
}

func getSecondCPM(floorCPM int64, ef []entity.Advertise) int64 {
	var secondCPM = floorCPM
	if len(ef) > 1 {
		secondCPM = ef[1].CPM()
	}

	return secondCPM
}

// winnerBid calculate winner bid
func winnerBid(cpm int64, ctr float64) int64 {
	return int64(float64(cpm)/(ctr*10)) + 1
}
