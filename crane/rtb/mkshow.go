package rtb

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"time"

	"clickyab.com/crane/crane/builder"
	"clickyab.com/crane/crane/capping"
	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/framework/router"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/safe"
)

func getSecondCPM(floorCPM int64, exceedFloor []entity.Advertise) int64 {
	var secondCPM = floorCPM
	if len(exceedFloor) > 1 && exceedFloor[0].Capping().Selected() == exceedFloor[1].Capping().Selected() {
		secondCPM = exceedFloor[1].CPM()
	}

	return secondCPM
}

func doBid(adData entity.Advertise, website entity.Publisher, slot entity.Seat, floorDiv int64) bool {
	adData.SetCTR((adData.AdCTR()*float64(adCTREffect.Int()) + slot.SlotCTR()*float64(slotCTREffect.Int())) / float64(100))
	adData.SetCPM(int64(float64(adData.Campaign().MaxBID()) * adData.CTR() * 10.0))
	//exceed cpm floor
	if floorDiv < 1 {
		floorDiv = 1
	}
	return adData.CPM() >= website.FloorCPM()/floorDiv
}

// winnerBid calculate winner bid
func winnerBid(cpm int64, ctr float64) int64 {
	return int64(float64(cpm)/(ctr*10)) + 1
}

func createInnerLinks(ctx *builder.Context) map[string]string {
	show := make(map[string]string)
	for _, slot := range ctx.GetRTB().Slots {
		u := url.URL{
			Scheme: ctx.GetCommon().Scheme,
			Host:   ctx.GetCommon().Host,
			Path: router.MustPath("show", map[string]string{
				"typ":     ctx.GetCommon().Type,
				"wid":     fmt.Sprint(ctx.GetPublisher().ID()),
				"mega":    ctx.GetCommon().MegaImp,
				"reserve": slot.ReservedHash(),
			}),
		}
		v := url.Values{}
		v.Set("tid", ctx.GetCommon().TID)
		v.Set("ref", ctx.GetCommon().Referrer)
		v.Set("loc", ctx.GetCommon().Parent)
		v.Set("s", fmt.Sprint(slot.ID))

		for i, j := range slot.ExtraParams() {
			v.Set(i, j)
		}
		u.RawQuery = v.Encode()
		show[slot.PublicID()] = u.String()
		slot.SetShowURL(u.String())
	}

	return show
}

// WARNING : DO NOT ADD PARAMETER TO THIS FUNCTION
func internalSelect(
	ctx *builder.Context,
	ads map[int][]entity.Advertise,
	wait chan map[string]entity.Advertise,
	winnerAd map[string]entity.Advertise,
	eav kv.Kiwi,
) {
	var noVideo bool // once set, never unset it again
	resAds := make(map[string]entity.Advertise)
	defer func() {
		if !ctx.GetRTB().Async {
			wait <- resAds
		}
	}()

	// TODO : must loop over this values, from lowest data to highest. the size with less ad count must be in higher priority
	selected := make(map[int]int)
	total := make(map[int]int)

	for _, slot := range ctx.GetRTB().Slots {
		var (
			exceedFloor []entity.Advertise
			underFloor  []entity.Advertise
		)

		for _, adData := range ads[slot.Size()] {
			total[slot.Size()]++
			if adData.Type() == entity.AdTypeVideo && noVideo {
				continue
			}
			if adData.WinnerBID() == 0 && doBid(adData, ctx.GetPublisher(), slot, ctx.GetRTB().FloorDIV) {
				exceedFloor = append(exceedFloor, adData)
			} else if adData.WinnerBID() == 0 {
				underFloor = append(underFloor, adData)
			}
		}
		var sorted []entity.Advertise
		var (
			ef     byMulti
			secBid bool
		)

		// order is to get data from exceed flor, then capp passed and if the config allowed,
		// use the under floor. for under floor there is no second biding pricing
		if len(exceedFloor) > 0 {
			ef = byMulti{
				Ads:   exceedFloor,
				Video: ctx.GetRTB().MultiVideo,
			}
			secBid = true
		} else if ctx.GetRTB().UnderFloor && len(underFloor) > 0 {
			ef = byMulti{
				Ads:   underFloor,
				Video: ctx.GetRTB().MultiVideo,
			}
			secBid = false
		}
		s := kv.GetSyncStore()
		if len(ef.Ads) == 0 {
			resAds[slot.PublicID()] = nil
			s.Push(slot.ReservedHash(), "no add", time.Hour)
			continue
		}

		sort.Sort(ef)
		sorted = ef.Ads

		theAd := sorted[0].Duplicate()

		// Do not do second biding pricing on this ads, they can not pass CPMFloor
		if secBid {
			secondCPM := getSecondCPM(ctx.GetPublisher().FloorCPM(), sorted)
			theAd.SetWinnerBID(winnerBid(secondCPM, theAd.CTR()), true)
		} else {
			theAd.SetWinnerBID(theAd.Campaign().MaxBID(), true)
		}

		if theAd.WinnerBID() > theAd.Campaign().MaxBID() {
			// TODO : must not happen, but it happen some how. check it later
			// since we change the winner bid, do not inc the cap
			theAd.SetWinnerBID(theAd.Campaign().MaxBID(), false)
		}

		if theAd.WinnerBID() < ctx.GetRTB().MinCPC {
			// since we change the winner bid, do not inc the cap
			theAd.SetWinnerBID(ctx.GetRTB().MinCPC, false)
		}

		theAd.SetSlot(slot)

		winnerAd[slot.PublicID()] = theAd
		resAds[slot.PublicID()] = theAd

		if !ctx.GetRTB().MultiVideo {
			noVideo = noVideo || theAd.Type() == entity.AdTypeVideo
		}
		eav.SetSubKey(
			fmt.Sprintf("AD_%d", theAd.ID()), fmt.Sprint(theAd.WinnerBID()),
		).SetSubKey(
			fmt.Sprintf("S_%d", theAd.ID()), fmt.Sprintf(slot.PublicID()),
		)
		s.Push(slot.ReservedHash(), fmt.Sprintf("%d", theAd.ID()), time.Hour)
		assert.Nil(capping.StoreCapping(ctx.GetCommon().CopID, theAd.ID()))
		selected[slot.Size()]++
	}
}

// selectAds is the only function that one must call to get ads
func selectAds(_ context.Context, ctx *builder.Context, ads map[int][]entity.Advertise) (map[string]string, map[string]entity.Advertise) {
	var (
		winnerAd = make(map[string]entity.Advertise)
	)

	links := createInnerLinks(ctx)

	var wait chan map[string]entity.Advertise
	if !ctx.GetRTB().Async {
		wait = make(chan map[string]entity.Advertise)
	}

	eav := kv.NewEavStore("MGA_" + ctx.GetCommon().MegaImp)
	err := eav.SetSubKey(
		"ip", ctx.GetCommon().IP.String(),
	).SetSubKey(
		"UA", ctx.GetCommon().UserAgent,
	).SetSubKey(
		"WS", fmt.Sprint(ctx.GetPublisher().ID()),
	).SetSubKey(
		"T", fmt.Sprint(time.Now().Unix()),
	).Save(megaImpExpire.Duration())
	assert.Nil(err)

	if !ctx.GetRTB().NoCap {
		ep := ctx.GetRTB().EventPage
		ads = capping.GetCapping(ctx.GetCommon().CopID, ads, ep, ctx.GetRTB().Slots...)
	} else {
		ads = capping.EmptyCapping(ads)
	}
	safe.GoRoutine(func() {
		internalSelect(ctx, ads, wait, winnerAd, eav)
	})
	var selected map[string]entity.Advertise
	if !ctx.GetRTB().Async {
		selected = <-wait
	}
	return links, selected
}
