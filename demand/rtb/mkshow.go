package rtb

import (
	"context"
	"fmt"
	"sort"

	"clickyab.com/crane/models/item"

	"github.com/clickyab/services/config"

	"clickyab.com/crane/demand/capping"
	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/assert"
)

var forceFristBid = config.RegisterBoolean("crane.demand.select.force_first_bid", true, "if it's set we ignore second bid")

func getSecondCPM(floorCPM float64, exceedFloor []entity.SelectedCreative) float64 {

	if forceFristBid.Bool() || !exceedFloor[0].IsSecBid() {
		return exceedFloor[0].CalculatedCPM()
	}

	if len(exceedFloor) > 1 &&
		exceedFloor[1].IsSecBid() &&
		!exceedFloor[1].Capping().Selected() &&
		exceedFloor[1].CalculatedCPM()+10 <= exceedFloor[0].CalculatedCPM() {
		return exceedFloor[1].CalculatedCPM() + 10
	}

	return floorCPM
}

func defaultCTR(seatType entity.RequestType, pub entity.PublisherType, sup entity.Supplier) float32 {
	return sup.DefaultCTR(fmt.Sprint(seatType), fmt.Sprint(pub))
}

func doBid(ad entity.Creative, slot entity.Seat, minCPM, minCPC float64, pub entity.Publisher) (float64, float64, float64, bool) {
	slotCtr := slot.CTR()
	if slot.CTR() < 0 {
		// get ctr based on the creative and seat type native app / native web / vast web ...
		slotCtr = defaultCTR(slot.RequestType(), pub.Type(), pub.Supplier())
	}
	adCtr := ad.AdCTR()
	if adCtr < 0 {
		// get ctr based on the creative and seat type native app / native web / vast web ...
		adCtr = defaultCTR(slot.RequestType(), pub.Type(), pub.Supplier())
	}
	ctr := float64(adCtr*float32(adCTREffect.Int())+slotCtr*float32(slotCTREffect.Int())) / float64(100)
	var cpc, cpm float64
	var exceed bool
	if ad.Campaign().Strategy() == entity.StrategyCPC {
		cpm = float64(ad.Campaign().MaxBID()) * ctr * 10.0
		cpc = float64(ad.Campaign().MaxBID())
		exceed = cpc >= minCPC
	} else {
		cpm = float64(ad.Campaign().MaxBID())
		cpc = float64(ad.Campaign().MaxBID()) / (ctr * 10.0)
		exceed = cpm >= minCPM
	}

	return ctr, cpm, cpc, exceed
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
	cps []entity.Campaign,
) {
	var noVideo bool                 // once set, never unset it again
	selected := make(map[int32]bool) // all ad selected in this session, to make sure they are not repeated

	for _, seat := range ctx.Seats() {
		ads := make([]entity.Creative, 0)
		for e := range cps {
			if ak, ok := cps[e].Sizes()[seat.Size()]; ok {
				ads = append(ads, ak...)
			}
		}

		if seat.RequestType() == entity.RequestTypeNative {
			ads = append(ads, target(ctx.User(), seat, cps)...)
		}

		exceedFloor, underFloor := selector(ctx, ads, seat, noVideo, selected)
		var (
			sorted []entity.SelectedCreative
			ef     byMulti
		)

		if len(exceedFloor) > 0 {
			ef = byMulti{
				Ads:   exceedFloor,
				Video: ctx.MultiVideo(),
			}
		} else if ctx.UnderFloor() && len(underFloor) > 0 {
			// under floor means we want to fill the seat at any cost. normally our own seat
			ef = byMulti{
				Ads:   underFloor,
				Video: ctx.MultiVideo(),
			}
		} else {
			continue
		}

		ef.Ads = capping.ApplyCapping(ctx.Capping(), ctx.User().ID(), ef.Ads)
		if len(ef.Ads) == 0 {
			continue
		}

		sort.Sort(ef)
		sorted = ef.Ads

		theAd := sorted[0]
		// Do not do second biding pricing on this ads, they can not pass CPMFloor
		targetCPM := getSecondCPM(seat.SoftCPM(), sorted)
		targetCPC := targetCPM / (theAd.CalculatedCTR() * 10.0)
		targetCPC, targetCPM = fixPrice(theAd.Campaign().Strategy(), targetCPC, targetCPM, seat.MinCPC(), seat.MinCPM())

		if targetCPM > float64(theAd.Campaign().MaxBID()) {
			targetCPM = float64(theAd.Campaign().MaxBID())
		}

		selected[theAd.ID()] = true

		// Only decrease share for CPM (which is reported to supplier) not bid (which is used by us)
		seat.SetWinnerAdvertise(theAd, targetCPC, targetCPM)

		if !ctx.MultiVideo() {
			noVideo = noVideo || theAd.Type() == entity.AdTypeVideo
		}
	}
}

func target(u entity.User, s entity.Seat, c []entity.Campaign) []entity.Creative {
	crs := make([]entity.Creative, 0)

	for e := range c {
		for _, v := range c[e].ReTargeting() {
			if ls, ok := u.List()[v]; ok {
				for i := range ls {
					it := item.GetItem(context.Background(), ls[i])
					if it == nil {
						continue
					}
					it.SetCampaign(c[e])
					crs = append(crs, it)
				}
			}
		}
	}
	return crs
}

func fixPrice(strategy entity.Strategy, cpc, cpm, minCPC, minCPM float64) (float64, float64) {

	if strategy == entity.StrategyCPC && cpc < minCPC {
		return minCPC, cpm
	}
	if strategy == entity.StrategyCPM && cpm < minCPM {
		return cpc, minCPM
	}
	return cpc, cpm
}

// selectAds is the only function that one must call to get ads
func selectAds(_ context.Context, ctx entity.Context, ads []entity.Campaign) {
	internalSelect(ctx, ads)
}

func selector(ctx entity.Context, ads []entity.Creative, seat entity.Seat, noVideo bool, selected map[int32]bool) (exceedFloor []entity.SelectedCreative, underFloor []entity.SelectedCreative) {
	assert.True(seat.SoftCPM() >= seat.MinCPM())

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

		if ctr, cpm, cpc, exceed := doBid(creative, seat, seat.MinCPM(), seat.MinCPC(), ctx.Publisher()); exceed {
			// a pass!
			exceedFloor = append(
				exceedFloor,
				adAndBid{
					Creative: creative,
					ctr:      ctr,
					cpm:      cpm,
					secBid:   cpm >= seat.SoftCPM(),
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
