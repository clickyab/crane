package capping

import (
	"fmt"
	"strconv"
	"time"

	"strings"

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
func noCappingMode(ads []entity.SelectedCreative) []entity.SelectedCreative {
	c := newContext()
	for i := range ads {
		capp := NewCapping(
			c,
			ads[i].Campaign().ID(),
			ads[i].Campaign().Frequency(),
			entity.CappingNone,
			"",
		)
		ads[i].SetCapping(capp)

	}
	return ads
}

// ApplyCapping is an entry for set capping in ads
func ApplyCapping(mode entity.CappingMode, copID string, ads []entity.SelectedCreative, ep string) []entity.SelectedCreative {
	switch mode {
	case entity.CappingNone:
		return noCappingMode(ads)
	case entity.CappingReset:
		return newResetCappingMode(copID, ads, ep)
		// return resetCappingMode(copID, ads, ep)
	case entity.CappingStrict:
		return newStrictCappingMode(copID, ads, ep)
		// return strictCappingMode(copID, ads, ep)
	}
	panic("invalid capping mode")
}

func round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}

func calculateAdsFrequency(ads []entity.SelectedCreative) map[int64]int {
	frqs := make(map[int64]int, len(ads))
	totalCTR := 0.0

	for i := range ads {
		totalCTR += ads[i].CalculatedCTR()
	}

	totalCap := float64(len(ads) * 3)
	for i := range ads {
		perc := ads[i].CalculatedCTR() * 100 / totalCTR
		frqs[ads[i].ID()] = int(round(perc * totalCap / 100))
	}

	return frqs
}

func newStrictCappingMode(copID string, ads []entity.SelectedCreative, ep string) []entity.SelectedCreative {
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

	adsFrequency := calculateAdsFrequency(ads)
	resp := make([]entity.SelectedCreative, 0, len(ads))
	for i := range ads {
		key := fmt.Sprintf(
			"%s_%d",
			adKey,
			ads[i].ID(),
		)

		view := results[key]
		n := 2.0
		if adsFrequency[ads[i].ID()] > 0 {
			n = float64(view) / float64(adsFrequency[ads[i].ID()])
		}

		if n <= 1 && !selected[ads[i].ID()] {
			passed := ads[i]
			capp := NewCapping(
				c,
				ads[i].Campaign().ID(),
				adsFrequency[ads[i].ID()],
				entity.CappingStrict,
				copID,
			)
			capp.IncView(ads[i].ID(), int(view), false)
			passed.SetCapping(capp)
			resp = append(resp, passed)
		}
	}

	return resp
}

/* func strictCappingMode(copID string, ads []entity.SelectedCreative, ep string) []entity.SelectedCreative {
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

	resp := make([]entity.SelectedCreative, 0, len(ads))
	for i := range ads {
		key := fmt.Sprintf(
			"%s_%d",
			adKey,
			ads[i].ID(),
		)

		view := results[key]
		n := float64(view) / float64(ads[i].Campaign().Frequency())
		if n <= 1 && !selected[ads[i].ID()] {
			passed := ads[i]
			capp := NewCapping(
				c,
				ads[i].Campaign().ID(),
				ads[i].Campaign().Frequency(),
				entity.CappingStrict,
				copID,
			)
			capp.IncView(ads[i].ID(), int(view), false)
			passed.SetCapping(capp)
			resp = append(resp, passed)
		}
	}

	return resp
} */

// GetCapping try to get capping for current ad
func newResetCappingMode(copID string, ads []entity.SelectedCreative, ep string) []entity.SelectedCreative {
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
	has := make(map[int]int)
	done := make(map[int][]struct {
		Key  string
		View int64
		entity.SelectedCreative
	})

	adsFrequency := calculateAdsFrequency(ads)
	resp := make([]entity.SelectedCreative, 0, len(ads))
	for i := range ads {
		size := ads[i].Size()
		if _, ok := has[size]; !ok {
			has[size] = 0
		}

		key := fmt.Sprintf(
			"%s_%d",
			adKey,
			ads[i].ID(),
		)

		view := results[key]
		n := 2.0
		if adsFrequency[ads[i].ID()] > 0 {
			n = float64(view) / float64(adsFrequency[ads[i].ID()])
		}

		if n < 1 && !selected[ads[i].ID()] {
			capp := NewCapping(
				c,
				ads[i].Campaign().ID(),
				adsFrequency[ads[i].ID()],
				entity.CappingReset,
				copID,
			)
			capp.IncView(ads[i].ID(), int(view), false)
			ads[i].SetCapping(capp)
			resp = append(resp, ads[i])
			has[size] = has[size] + 1
		} else if !selected[ads[i].ID()] {
			// capping is passed
			done[size] = append(done[size], struct {
				Key  string
				View int64
				entity.SelectedCreative
			}{
				Key:              key,
				View:             view,
				SelectedCreative: ads[i],
			})
		}
	}

	for size := range has {
		// reset this size
		var sizedCap = make([]string, len(done[size]))
		for i := range done[size] {
			if has[size] == 0 {
				done[size][i].View = 0
			}
			capp := NewCapping(
				c,
				done[size][i].Campaign().ID(),
				adsFrequency[done[size][i].ID()],
				entity.CappingReset,
				copID,
			)
			capp.IncView(done[size][i].ID(), int(done[size][i].View), false)
			sizedCap[i] = done[size][i].Key
			done[size][i].SetCapping(capp)
			resp = append(resp, done[size][i].SelectedCreative)
		}
		if has[size] == 0 {
			logrus.Debugf("remove key for size %d keys are %s", size, strings.Join(sizedCap, ","))
			_ = ck.Drop(sizedCap...)
		}
	}
	return resp
}

// GetCapping try to get capping for current ad
/* func resetCappingMode(copID string, ads []entity.SelectedCreative, ep string) []entity.SelectedCreative {
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
	has := make(map[int]int)
	done := make(map[int][]struct {
		Key  string
		View int64
		entity.SelectedCreative
	})
	resp := make([]entity.SelectedCreative, 0, len(ads))
	for i := range ads {
		size := ads[i].Size()
		if _, ok := has[size]; !ok {
			has[size] = 0
		}

		key := fmt.Sprintf(
			"%s_%d",
			adKey,
			ads[i].ID(),
		)

		view := results[key]
		n := float64(view) / float64(ads[i].Campaign().Frequency())
		if n < 1 && !selected[ads[i].ID()] {
			capp := NewCapping(
				c,
				ads[i].Campaign().ID(),
				ads[i].Campaign().Frequency(),
				entity.CappingReset,
				copID,
			)
			capp.IncView(ads[i].ID(), int(view), false)
			ads[i].SetCapping(capp)
			resp = append(resp, ads[i])
			has[size] = has[size] + 1
		} else if !selected[ads[i].ID()] {
			// capping is passed
			done[size] = append(done[size], struct {
				Key  string
				View int64
				entity.SelectedCreative
			}{
				Key:              key,
				View:             view,
				SelectedCreative: ads[i],
			})
		}
	}

	for size := range has {
		// reset this size
		var sizedCap = make([]string, len(done[size]))
		for i := range done[size] {
			if has[size] == 0 {
				done[size][i].View = 0
			}
			capp := NewCapping(
				c,
				done[size][i].Campaign().ID(),
				done[size][i].Campaign().Frequency(),
				entity.CappingReset,
				copID,
			)
			capp.IncView(done[size][i].ID(), int(done[size][i].View), false)
			sizedCap[i] = done[size][i].Key
			done[size][i].SetCapping(capp)
			resp = append(resp, done[size][i].SelectedCreative)
		}
		if has[size] == 0 {
			logrus.Debugf("remove key for size %d keys are %s", size, strings.Join(sizedCap, ","))
			_ = ck.Drop(sizedCap...)
		}
	}
	return resp
} */
