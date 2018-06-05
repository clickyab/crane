package distribution

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/clickyab/services/array"

	"clickyab.com/crane/demand/entity"

	"github.com/clickyab/services/assert"

	"github.com/clickyab/services/config"
)

var (
	ctrSoftFloorEffect       = config.RegisterInt("crane.distribution.ctr.softfloor.effect", 50, "percent of effect ctr soft floor. default is 50%")
	ctrSoftFloorReduceMargin = config.RegisterFloat64("crane.distribution.ctr.softfloor.reduce.margin", 20, "percent of reduce margin of ctr soft floor. default is 20%")
	upperCreatives           = make([]string, 0)
	underCreatives           = make([]string, 0)
	lock                     = &sync.Mutex{}
	firstTry                 = true
	ctrSoftFloor             = 0.039
	seatCnt                  = 1
	sortableCreatives        entity.SortableCreative
)

//GetSortStrategy get creatives sort strategy base on distribution ctr soft floor role and status
func GetSortStrategy(ctx entity.Context, seat entity.Seat) entity.SortStrategyType {
	seatCnt = len(ctx.Seats())

	seatID := seat.PublicID()
	if seatCnt > 1 {
		seatID = fmt.Sprintf("multi_%d", seatCnt)
	}
	setDistributeKey(ctx.Publisher().Supplier().Name(), seat.Type().String(), ctx.Publisher().ID(), seat.Size(), seatID)
	loadSelected()

	if len(upperCreatives) < getNeedUpCount() && len(underCreatives) == 0 {
		return entity.SortByCTR
	}

	return entity.SortByCPM
}

func getNeedUpCount() int {
	if seatCnt > 1 {
		return seatCnt * ctrSoftFloorEffect.Int() / 100
	}

	return 10 * ctrSoftFloorEffect.Int() / 100
}

func getTotalCount() int {
	if seatCnt > 1 {
		return seatCnt
	}

	return 10
}

//GetWinner to apply ctr soft floor ditribution role and return winner creative
func GetWinner(ctx entity.Context, sc entity.SortableCreative) entity.SelectedCreative {
	assert.NotNil(ctx.Seats())
	defer checkAndReset()

	sortableCreatives = sc
	sorted := sortableCreatives.Ads

	statistics := ctx.GetCreativesStatistics()
	for _, v := range statistics {
		if entity.AdType(v.CreativeType()) == sorted[0].Type() {
			ctrSoftFloor = (100 - ctrSoftFloorReduceMargin.Float64()) * v.AvgCTR() / 100
			break
		}
	}
	loadSelected()

	if seatCnt > 1 {
		return findWinnerMultiSeats(sorted)
	}

	return findWinnerSingleSeats(sorted)
}

func findWinnerMultiSeats(sorted []entity.SelectedCreative) entity.SelectedCreative {
	if len(upperCreatives) < getNeedUpCount() {
		for i := range sorted {
			crID := fmt.Sprintf("%d", sorted[i].ID())

			if !array.StringInArray(crID, upperCreatives...) && sorted[i].CalculatedCTR() >= ctrSoftFloor {
				selectUpper(crID)
				return sorted[i]
			}
		}
	}

	sorted = sortAgainByCPM()
	for i := range sorted {
		crID := fmt.Sprintf("%d", sorted[i].ID())

		if !array.StringInArray(crID, upperCreatives...) { //We don't check creative is in upper or not beacause this check will be break capping
			selectUnder(crID)
			return sorted[i]
		}
	}

	selectUnder(fmt.Sprintf("%d", sorted[0].ID()))
	return sorted[0]
}

func findWinnerSingleSeats(sorted []entity.SelectedCreative) entity.SelectedCreative {
	if len(upperCreatives) < getNeedUpCount() {
		for i := range sorted {
			if sorted[i].CalculatedCTR() >= ctrSoftFloor { //We don't check creative is in upper or not beacause this check will be break capping
				selectUpper(fmt.Sprintf("%d", sorted[i].ID()))
				return sorted[i]
			}
		}
	}

	sorted = sortAgainByCPM()

	//We don't check upper and under here. beacause don't want to break capping roles
	selectUnder(fmt.Sprintf("%d", sorted[0].ID()))
	return sorted[0]
}

func sortAgainByCPM() []entity.SelectedCreative {
	sortableCreatives.SortStrategy = entity.SortByCPM
	sort.Sort(sortableCreatives)

	return sortableCreatives.Ads
}

func loadSelected() {
	lock.Lock()
	defer lock.Unlock()
	selected := getDistributionSelected()

	if selected["upper"] != "" {
		upperCreatives = strings.Split(selected["upper"], ",")
	}
	if selected["under"] != "" {
		underCreatives = strings.Split(selected["under"], ",")
	}

	checkAndReset()
}

func checkAndReset() {
	if len(upperCreatives)+len(underCreatives) >= getTotalCount() || (seatCnt > 1 && firstTry) {
		firstTry = false
		upperCreatives = make([]string, 0)
		underCreatives = make([]string, 0)

		assert.Nil(setUpperCreative(""))
		assert.Nil(setUnderCreative(""))
	}
}

func selectUpper(id string) {
	lock.Lock()
	defer lock.Unlock()

	upperCreatives = append(upperCreatives, id)
	assert.Nil(setUpperCreative(strings.Join(upperCreatives, ",")))
}

func selectUnder(id string) {
	lock.Lock()
	defer lock.Unlock()

	underCreatives = append(underCreatives, id)
	assert.Nil(setUnderCreative(strings.Join(underCreatives, ",")))
}
