package rtb

import (
	"context"
	"sort"

	"fmt"

	"clickyab.com/crane/demand/capping"
	"clickyab.com/crane/demand/entity"
)

func getSecondCPM(floorCPM float64, exceedFloor []entity.SelectedCreative) float64 {
	if !exceedFloor[0].IsSecBid() {
		return float64(exceedFloor[0].CalculatedCPM())
	}

	if len(exceedFloor) > 1 &&
		exceedFloor[1].IsSecBid() &&
		!exceedFloor[1].Capping().Selected() &&
		exceedFloor[1].CalculatedCPM()+10 <= exceedFloor[0].CalculatedCPM() {
		return float64(exceedFloor[1].CalculatedCPM() + 10)
	}

	return floorCPM
}

func defaultCTR(seatType entity.RequestType, pub entity.PublisherType, sup entity.Supplier) float64 {
	return sup.DefaultCTR(fmt.Sprint(seatType), fmt.Sprint(pub))
}

func doBid(ad entity.Creative, slot entity.Seat, minCPM, minCPC float64, pub entity.Publisher) (float64, float64, float64, bool) {
	slotCtr := slot.CTR()
	if slot.CTR() < 0 {
		//get ctr based on the creative and seat type native app / native web / vast web ...
		slotCtr = defaultCTR(slot.RequestType(), pub.Type(), pub.Supplier())
	}
	adCtr := ad.AdCTR()
	if adCtr < 0 {
		//get ctr based on the creative and seat type native app / native web / vast web ...
		adCtr = defaultCTR(slot.RequestType(), pub.Type(), pub.Supplier())
	}
	ctr := (adCtr*float64(adCTREffect.Int()) + slotCtr*float64(slotCTREffect.Int())) / float64(100)
	var cpc, cpm float64
	var under bool
	if ad.Campaign().Strategy() == entity.StrategyCPC {
		cpm = float64(ad.MaxBID()) * ctr * 10.0
		cpc = float64(ad.MaxBID())
		under = float64(cpc) > minCPC
	} else {
		cpm = float64(ad.MaxBID())
		cpc = float64(ad.MaxBID()) / (ctr * 10.0)
		under = float64(cpm) > minCPM
	}

	return ctr, cpm, cpc, under
}

func incShare(sup entity.Supplier, price int64) int64 {
	return (price * int64(sup.Share())) / 100
}

func decShare(sup entity.Supplier, price float64) float64 {
	return (price * 100.0) / float64(sup.Share())
}

type adAndBid struct {
	entity.Creative
	ctr    float64
	cpm    float64
	cpc    float64
	secBid bool
}

func (aab adAndBid) CalculatedCPC() float64 {
	return aab.cpc
}

func (aab adAndBid) CalculatedCTR() float64 {
	return aab.ctr
}

func (aab adAndBid) CalculatedCPM() float64 {
	return aab.cpm
}

func (aab adAndBid) IsSecBid() bool {
	return aab.secBid
}

// WARNING : DO NOT ADD PARAMETER TO THIS FUNCTION
func internalSelect(
	ctx entity.Context,
	ads []entity.Creative,
) {
	var noVideo bool                 // once set, never unset it again
	selected := make(map[int64]bool) // all ad selected in this session, to make sure they are not repeated

	for _, seat := range ctx.Seats() {
		var (
			softCPM = float64(ctx.Publisher().Supplier().SoftFloorCPM(fmt.Sprint(seat.RequestType()), fmt.Sprint(ctx.Publisher().Type()))) // softCPM floor , determine the sec bidding pricing
			minCPM  = float64(incShare(ctx.Publisher().Supplier(), seat.MinBid()))                                                         // minimum cpm of this seat, aka hard floor, after adding our share to it
			minCPC  = float64(ctx.Publisher().Supplier().SoftFloorCPC(fmt.Sprint(seat.RequestType()), fmt.Sprint(ctx.Publisher().Type()))) // minimum cpc of this seat, aka hard floor, after adding our share to it
		)

		// softCPM floor is smaller than hard floor, so we do not have sec biding
		if softCPM < minCPM {
			softCPM = minCPM
		}

		exceedFloor, underFloor := selector(ctx, ads, seat, noVideo, selected)

		var (
			sorted []entity.SelectedCreative
			ef     byMulti
		)

		// order is to get data from exceed floor, then capp passed and if the config allowed,
		// use the under floor. for under floor there is no second biding pricing
		var under bool
		if len(exceedFloor) > 0 {
			ef = byMulti{
				Ads:   exceedFloor,
				Video: ctx.MultiVideo(),
			}
		} else if ctx.UnderFloor() && len(underFloor) > 0 {
			// under floor means we want to fill the seat at any cost. normally our own seat
			under = true
			ef = byMulti{
				Ads:   underFloor,
				Video: ctx.MultiVideo(),
			}
		} else {
			continue
		}

		ef.Ads = capping.ApplyCapping(ctx.Capping(), ctx.User().ID(), ef.Ads, ctx.EventPage())
		if len(ef.Ads) == 0 {
			continue
		}
		sort.Sort(ef)
		sorted = ef.Ads

		theAd := sorted[0]
		// Do not do second biding pricing on this ads, they can not pass CPMFloor
		targetCPM := getSecondCPM(softCPM, sorted)
		targetCPC := targetCPM / (theAd.CalculatedCTR() * 10.0)

		if !under {
			targetCPC, targetCPM = fixPrice(theAd.Campaign().Strategy(), targetCPC, targetCPM, minCPC, minCPM)
		}

		selected[theAd.ID()] = true
		// Only decrease share for CPM (which is reported to supplier) not bid (which is used by us)
		seat.SetWinnerAdvertise(theAd, targetCPC, decShare(ctx.Publisher().Supplier(), targetCPM))

		if !ctx.MultiVideo() {
			noVideo = noVideo || theAd.Type() == entity.AdTypeVideo
		}
		// TODO : The real problem is what if we are not going to win? this assume any select means show.
		theAd.Capping().Store(theAd.ID())
	}
}

func fixPrice(strategy entity.Strategy, cpc, cpm, minCPC, minCPM float64) (float64, float64) {
	if strategy == entity.StrategyCPC && cpm < minCPM {
		return cpc, minCPM
	}
	if strategy == entity.StrategyCPM && cpc < minCPC {
		return minCPC, cpm
	}
	return cpc, cpm
}

// selectAds is the only function that one must call to get ads
func selectAds(_ context.Context, ctx entity.Context, ads []entity.Creative) {
	internalSelect(ctx, ads)
}

func selector(ctx entity.Context, ads []entity.Creative, seat entity.Seat, noVideo bool, selected map[int64]bool) ([]entity.SelectedCreative, []entity.SelectedCreative) {
	var (
		exceedFloor []entity.SelectedCreative // above  hard floor (the minimum cpm ), legit ads
		underFloor  []entity.SelectedCreative
		// not passed from floor, only used if the supplier accept less than minCPM bids, normally only us, as clickyab
		softCPM = ctx.Publisher().Supplier().SoftFloorCPM(fmt.Sprint(seat.RequestType()), fmt.Sprint(ctx.Publisher().Type()))          // softCPM floor , determine the sec bidding pricing
		minCPM  = incShare(ctx.Publisher().Supplier(), seat.MinBid())                                                                  // minimum cpm of this seat, aka hard floor, after adding our share to it
		minCPC  = float64(ctx.Publisher().Supplier().SoftFloorCPC(fmt.Sprint(seat.RequestType()), fmt.Sprint(ctx.Publisher().Type()))) // minimum cpc of this seat, aka hard floor, after adding our share to it

	)

	// softCPM floor is smaller than hard floor, so we do not have sec biding
	if softCPM < minCPM {
		softCPM = minCPM
	}

	for _, creative := range ads {
		if creative.Type() == entity.AdTypeVideo && noVideo {
			continue
		}
		if selected[creative.ID()] {
			continue
		}
		if !seat.Acceptable(creative) {
			continue
		}

		if ctr, cpm, cpc, ok := doBid(creative, seat, float64(minCPM), minCPC, ctx.Publisher()); ok {
			// a pass!
			exceedFloor = append(
				exceedFloor,
				adAndBid{
					Creative: creative,
					ctr:      ctr,
					cpm:      cpm,
					secBid:   cpm >= float64(softCPM),
					cpc:      cpc,
				},
			)
		} else {
			underFloor = append(
				underFloor,
				adAndBid{
					Creative: creative,
					ctr:      ctr,
					cpm:      cpm,
					cpc:      cpc,
					secBid:   false,
				},
			)
		}
	}
	return exceedFloor, underFloor
}
