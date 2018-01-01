package capping

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/kv"
	"github.com/sirupsen/logrus"
)

const (
	cappingKey = "CP"
	adKey      = "AD"
)

var (
	dailyCapExpire = config.RegisterDuration("crane.capping.daily_cap_expire", 24*time.Hour, "daily capping expire")
)

func getCappingKey(mode entity.CappingMode, copID string) string {
	return fmt.Sprintf(
		"%s_%d_%s_%s",
		cappingKey,
		int(mode),
		copID,
		time.Now().Format("060102"),
	)
}

// EmptyCapping is a hack to handle no capping situation
func noCappingMode(ads map[int][]entity.Advertise) map[int][]entity.Advertise {
	c := newContext()
	for i := range ads {
		for j := range ads[i] {
			capp := NewCapping(
				c,
				ads[i][j].Campaign().ID(),
				0,
				ads[i][j].Campaign().Frequency(),
			)
			ads[i][j].SetCapping(capp)

		}
		sortCap := sortByCap(ads[i])
		sort.Sort(sortCap)
		ads[i] = []entity.Advertise(sortCap)
	}

	return ads
}

// ApplyCapping is an entry for set capping in ads
func ApplyCapping(mode entity.CappingMode, copID string, ads map[int][]entity.Advertise, ep string, slots ...entity.Seat) map[int][]entity.Advertise {
	switch mode {
	case entity.CappingNone:
		return noCappingMode(ads)
	case entity.CappingReset:
		return resetCappingMode(copID, ads, ep, slots...)
	case entity.CappingStrict:
		return strictCappingMode(copID, ads, ep, slots...)
	}
	panic("invalid capping mode")
}

func strictCappingMode(copID string, ads map[int][]entity.Advertise, ep string, slots ...entity.Seat) map[int][]entity.Advertise {
	var selected = make(map[int64]bool)

	// evenet page is an old hack to handle ads in same page in multiple request. maybe we should retire it
	// TODO : remove event page after 31 March 2018 if there is no need for it
	if ep != "" {
		s := kv.NewDistributedSet(ep)

		for _, v := range s.Members() {
			vInt, _ := strconv.ParseInt(v, 10, 0)
			selected[vInt] = true
		}

	}
	c := newContext()
	ck := kv.NewAEAVStore(getCappingKey(entity.CappingStrict, copID), dailyCapExpire.Duration())
	results := ck.AllKeys()
	doneSized := make(map[int]bool)
	resp := make(map[int][]entity.Advertise)
	for i := range slots {
		size := slots[i].Size()
		if doneSized[size] {
			continue
		}
		doneSized[size] = true
		for ad := range ads[size] {
			key := fmt.Sprintf(
				"%s_%d",
				adKey,
				ads[size][ad].ID(),
			)
			view := results[key]
			n := float64(view) / float64(ads[size][ad].Campaign().Frequency())
			if n <= 1 && !selected[ads[size][ad].ID()] {
				passed := ads[size][ad]
				capp := NewCapping(
					c,
					ads[size][ad].Campaign().ID(),
					0,
					ads[size][ad].Campaign().Frequency(),
				)
				capp.IncView(passed.ID(), int(view), false)
				passed.SetCapping(capp)
				resp[size] = append(resp[size], passed)
			}
		}
	}
	return resp
}

// GetCapping try to get capping for current ad
func resetCappingMode(copID string, ads map[int][]entity.Advertise, ep string, slots ...entity.Seat) map[int][]entity.Advertise {
	var selected = make(map[int64]bool)
	// evenet page is an old hack to handle ads in same page in multiple request. maybe we should retire it
	// TODO : remove event page after 31 March 2018 if there is no need for it
	if ep != "" {
		s := kv.NewDistributedSet(ep)

		for _, v := range s.Members() {
			vInt, _ := strconv.ParseInt(v, 10, 0)
			selected[vInt] = true
		}

	}
	c := newContext()
	ck := kv.NewAEAVStore(getCappingKey(entity.CappingReset, copID), dailyCapExpire.Duration())
	results := ck.AllKeys()
	doneSized := make(map[int]bool)
	for i := range slots {
		size := slots[i].Size()
		if doneSized[size] {
			continue
		}
		doneSized[size] = true
		found := false
		var sizeCap []string
		for ad := range ads[size] {
			key := fmt.Sprintf(
				"%s_%d",
				adKey,
				ads[size][ad].ID(),
			)
			view := results[key]
			sizeCap = append(sizeCap, key)
			n := float64(view) / float64(ads[size][ad].Campaign().Frequency())
			if n < 1 {
				found = true
				break // there is still one campaign
			}
		}
		// if not found then reset all capping
		if !found {
			logrus.Debugf("Removing key for size %d", size)
			_ = ck.Drop(sizeCap...)
			for i := range sizeCap {
				results[sizeCap[i]] = 0
			}
		}
		for ad := range ads[size] {
			var view int64
			if found {
				view = results[fmt.Sprintf(
					"%s_%d",
					adKey,
					ads[size][ad].ID(),
				)]
			}
			capp := NewCapping(
				c,
				ads[size][ad].Campaign().ID(),
				0,
				ads[size][ad].Campaign().Frequency(),
			)
			capp.IncView(ads[size][ad].ID(), int(view), selected[ads[size][ad].ID()])
			ads[size][ad].SetCapping(capp)
		}
	}
	return ads
}

// StoreCapping try to store a capping object
func StoreCapping(mode entity.CappingMode, copID string, adID int64) int64 {
	return kv.NewAEAVStore(getCappingKey(mode, copID), dailyCapExpire.Duration()).IncSubKey(fmt.Sprintf("%s_%d", adKey, adID), 1)
}
