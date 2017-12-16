package rtb

import (
	"context"
	"sort"

	"clickyab.com/crane/crane/capping"
	"clickyab.com/crane/crane/entity"
)

func getSecondCPM(floorCPM int64, exceedFloor []adAndBid) float64 {
	var secondCPM = float64(floorCPM)
	if len(exceedFloor) > 1 && exceedFloor[0].Capping().Selected() == exceedFloor[1].Capping().Selected() {
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

type adAndBid struct {
	entity.Advertise
	ctr float64
	cpm int64
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
			exceedFloor []adAndBid
			underFloor  []adAndBid
		)

		for _, adData := range ads[seat.Size()] {
			if adData.Type() == entity.AdTypeVideo && noVideo {
				continue
			}
			if selected[adData.ID()] {
				continue
			}
			// Drop if the min bid is lesser than requested min bid
			if adData.Campaign().MaxBID() < seat.MinBid() {
				continue
			}
			if ctr, cpm, ok := doBid(adData, seat, ctx.FloorCPM()); ok {
				exceedFloor = append(exceedFloor, adAndBid{Advertise: adData, ctr: ctr, cpm: cpm})
			} else {
				underFloor = append(underFloor, adAndBid{Advertise: adData, ctr: ctr, cpm: cpm})
			}
		}
		var sorted []adAndBid
		var (
			ef     byMulti
			secBid bool
		)

		// order is to get data from exceed flor, then capp passed and if the config allowed,
		// use the under floor. for under floor there is no second biding pricing
		if len(exceedFloor) > 0 {
			ef = byMulti{
				Ads:   exceedFloor,
				Video: ctx.MultiVideo(),
			}
			secBid = true
		} else if len(underFloor) > 0 {
			ef = byMulti{
				Ads:   underFloor,
				Video: ctx.MultiVideo(),
			}
			secBid = false
		} else {
			continue
		}

		sort.Sort(ef)
		sorted = ef.Ads

		theAd := sorted[0]
		// Do not do second biding pricing on this ads, they can not pass CPMFloor
		targetCPM := float64(theAd.Campaign().MaxBID())
		if secBid {
			targetCPM = getSecondCPM(ctx.FloorCPM(), sorted)
		}
		bid := winnerBid(targetCPM, theAd.ctr)
		if bid > float64(theAd.Campaign().MaxBID()) {
			// TODO : must not happen, but it happen some how. check it later
			// since we change the winner bid, do not inc the cap
			bid = float64(theAd.Campaign().MaxBID())
			// also fix target cpm
			targetCPM = theAd.ctr * 10 * bid
		}

		if bid < float64(seat.MinBid()) {
			// since we change the winner bid, do not inc the cap
			bid = float64(seat.MinBid())
			// also fix target cpm
			targetCPM = theAd.ctr * 10 * bid
		}
		selected[theAd.ID()] = true
		seat.SetWinnerAdvertise(theAd.Advertise, bid, targetCPM)

		if !ctx.MultiVideo() {
			noVideo = noVideo || theAd.Type() == entity.AdTypeVideo
		}
		capping.StoreCapping(
			ctx.User().ID(),
			theAd.ID())
	}
}

// selectAds is the only function that one must call to get ads
func selectAds(_ context.Context, ctx entity.Context, ads map[int][]entity.Advertise) {
	if ctx.Capping() {
		ep := ctx.EventPage()
		ads = capping.GetCapping(ctx.User().ID(), ads, ep, ctx.Seats()...)
	} else {
		ads = capping.EmptyCapping(ads)
	}
	internalSelect(ctx, ads)
}
