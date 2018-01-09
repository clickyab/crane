package rtb

import (
	"context"
	"sort"

	"clickyab.com/crane/demand/capping"
	"clickyab.com/crane/demand/entity"
)

func getSecondCPM(floorCPM int64, exceedFloor []adAndBid) float64 {
	var secondCPM = float64(floorCPM)
	if len(exceedFloor) > 1 && // if there is more than one
		exceedFloor[1].secBid && // the next is also a second bid ad
		exceedFloor[0].Capping().Selected() == exceedFloor[1].Capping().Selected() { // and second is not selected already
		secondCPM = float64(exceedFloor[1].cpm)
	}

	return secondCPM
}

func doBid(adData entity.Advertise, slot entity.Seat, floorCPM int64) (float64, int64, bool) {
	ctr := (adData.AdCTR()*float64(adCTREffect.Int()) + slot.CTR()*float64(slotCTREffect.Int())) / float64(100)
	cpm := int64(float64(adData.Campaign().MaxBID()) * ctr * 10.0)
	return ctr, cpm, cpm >= floorCPM
}

// winnerBid calculate winner bid
func winnerBid(cpm float64, ctr float64) float64 {
	return cpm / (ctr * 10)
}

func incShare(sup entity.Supplier, price int64) int64 {
	return (price * int64(sup.Share())) / 100
}

func decShare(sup entity.Supplier, price float64) float64 {
	return (price * 100.0) / float64(sup.Share())
}

type adAndBid struct {
	entity.Advertise
	ctr    float64
	cpm    int64
	secBid bool
}

// WARNING : DO NOT ADD PARAMETER TO THIS FUNCTION
func internalSelect(
	ctx entity.Context,
	ads map[int][]entity.Advertise,
) {
	var noVideo bool                 // once set, never unset it again
	selected := make(map[int64]bool) // all ad selected in this session, to make sure they are not repeated

	for _, seat := range ctx.Seats() {
		var (
			exceedFloor []adAndBid                                            // above  hard floor (the minimum cpm ), legit ads
			underFloor  []adAndBid                                            // not passed from floor, only used if the supplier accept less than minCPM bids, normally only us, as clickyab
			soft        = ctx.Publisher().SoftFloorCPM()                      // soft floor , determine the sec bidding pricing
			minCPM      = incShare(ctx.Publisher().Supplier(), seat.MinBid()) // minimum cpm of this seat, aka hard floor, after adding our share to it
		)

		// soft floor is smaller than hard floor, so we do not have sec biding
		if soft < minCPM {
			soft = minCPM
		}

		for _, adData := range ads[seat.Size()] {
			if !seat.Acceptable(adData) {
				continue
			}
			if adData.Type() == entity.AdTypeVideo && noVideo {
				continue
			}
			if selected[adData.ID()] {
				continue
			}

			if ctr, cpm, ok := doBid(adData, seat, minCPM); ok {
				// a pass!
				exceedFloor = append(
					exceedFloor,
					adAndBid{
						Advertise: adData,
						ctr:       ctr,
						cpm:       cpm,
						secBid:    cpm >= soft,
					},
				)
			} else {
				underFloor = append(
					underFloor,
					adAndBid{
						Advertise: adData,
						ctr:       ctr,
						cpm:       cpm,
						secBid:    false,
					},
				)
			}
		}

		var (
			sorted []adAndBid
			ef     byMulti
		)

		// order is to get data from exceed floor, then capp passed and if the config allowed,
		// use the under floor. for under floor there is no second biding pricing
		if len(exceedFloor) > 0 {
			ef = byMulti{
				Ads:   exceedFloor,
				Video: ctx.MultiVideo(),
			}
		} else if ctx.Publisher().Supplier().UnderFloor() && len(underFloor) > 0 {
			// under floor means we want to fill the seat at any cost. normally our own seat
			ef = byMulti{
				Ads:   underFloor,
				Video: ctx.MultiVideo(),
			}
		} else {
			continue
		}

		sort.Sort(ef)
		sorted = ef.Ads

		theAd := sorted[0]
		// Do not do second biding pricing on this ads, they can not pass CPMFloor
		targetCPM := float64(theAd.Campaign().MaxBID()) * 10 * theAd.ctr
		if theAd.secBid {
			targetCPM = getSecondCPM(ctx.SoftFloorCPM(), sorted)
		}

		// bid is in CPC world, so must compare it with the max bid
		bid := winnerBid(targetCPM, theAd.ctr)
		if bid > float64(theAd.Campaign().MaxBID()) {
			// TODO : must not happen, but it happen some how. check it later
			// since we change the winner bid, do not inc the cap
			bid = float64(theAd.Campaign().MaxBID())
			// also fix target cpm
			targetCPM = theAd.ctr * 10 * bid
		}

		// minCPM is in CPM world. so must compare it with target CPM
		if targetCPM < float64(minCPM) {
			targetCPM = float64(minCPM)
		}
		selected[theAd.ID()] = true
		// Only decrease share for CPM (which is reported to supplier) not bid (which is used by us)
		seat.SetWinnerAdvertise(theAd.Advertise, bid, decShare(ctx.Publisher().Supplier(), targetCPM))

		if !ctx.MultiVideo() {
			noVideo = noVideo || theAd.Type() == entity.AdTypeVideo
		}
		// TODO : The real problem is what if we are not going to win? this assume any select means show.
		capping.StoreCapping(
			ctx.Capping(),
			ctx.User().ID(),
			theAd.ID())
	}
}

// selectAds is the only function that one must call to get ads
func selectAds(_ context.Context, ctx entity.Context, ads map[int][]entity.Advertise) {
	ads = capping.ApplyCapping(ctx.Capping(), ctx.User().ID(), ads, ctx.EventPage(), ctx.Seats()...)
	internalSelect(ctx, ads)
}
