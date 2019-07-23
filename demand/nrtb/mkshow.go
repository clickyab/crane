package nrtb

import (
	"context"
	"fmt"
	"sort"

	"clickyab.com/crane/models/item"

	"clickyab.com/crane/demand/capping"
	"clickyab.com/crane/demand/entity"
)

func defaultCTR(seatType entity.RequestType, pub entity.PublisherType, sup entity.Supplier) float32 {
	return sup.DefaultCTR(fmt.Sprint(seatType), fmt.Sprint(pub))
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

// selectAds is the only function that one must call to get ads
// WARNING : DO NOT ADD PARAMETER TO THIS FUNCTION
func selectAds(
	ctx entity.Context,
	cps []entity.Campaign,
) {

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
		exceedFloor := make([]entity.SelectedCreative, 0)
		for _, creative := range ads {
			if selected[creative.ID()] {
				continue
			}
			if !seat.Acceptable(creative) {
				continue
			}
			slotCtr := seat.CTR()
			if seat.CTR() < 0 {
				// get ctr based on the creative and seat type native app / native web / vast web ...
				slotCtr = defaultCTR(seat.RequestType(), ctx.Publisher().Type(), ctx.Publisher().Supplier())
			}
			adCtr := creative.AdCTR()
			if adCtr < 0 {
				// get ctr based on the creative and seat type native app / native web / vast web ...
				adCtr = defaultCTR(seat.RequestType(), ctx.Publisher().Type(), ctx.Publisher().Supplier())
			}
			ctr := float64(adCtr*float32(adCTREffect.Int())+slotCtr*float32(slotCTREffect.Int())) / float64(100)
			var cpc, cpm float64
			var exceed bool
			if creative.Campaign().Strategy() == entity.StrategyCPC {
				cpm = float64(creative.Campaign().MaxBID()) * ctr * 10.0
				cpc = float64(creative.Campaign().MaxBID())
				exceed = cpc >= seat.MinCPC()
			} else {
				cpm = float64(creative.Campaign().MaxBID())
				cpc = float64(creative.Campaign().MaxBID()) / (ctr * 10.0)
				exceed = cpm >= seat.MinCPM()
			}

			if exceed {
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
			}
		}

		var (
			sorted []entity.SelectedCreative
			ef     byMulti
		)

		if len(exceedFloor) == 0 {
			continue
		}
		ef = byMulti{
			Ads:   exceedFloor,
			Video: ctx.MultiVideo(),
		}

		ef.Ads = capping.ApplyCapping(ctx.Capping(), ctx.User().ID(), ef.Ads)
		if len(ef.Ads) == 0 {
			continue
		}

		sort.Sort(ef)
		sorted = ef.Ads

		theAd := sorted[0]

		targetCPC := float64(theAd.Campaign().MaxBID())

		if ctx.Publisher().MaxCPC() > 0 && targetCPC > ctx.Publisher().MaxCPC() {
			targetCPC = ctx.Publisher().MaxCPC()
		}

		if float64(seat.MinBid()) > targetCPC {
			continue
		}

		selected[theAd.ID()] = true

		// Only decrease share for CPM (which is reported to supplier) not bid (which is used by us)
		seat.SetWinnerAdvertise(theAd, targetCPC, targetCPC)

	}
}
