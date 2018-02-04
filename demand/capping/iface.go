package capping

import (
	"fmt"
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
func noCappingMode(ads []entity.Creative) []entity.Creative {
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
	// No need to sort this one :)
	return ads
}

// ApplyCapping is an entry for set capping in ads
func ApplyCapping(mode entity.CappingMode, copID string, ads []entity.Creative, ep string) []entity.Creative {
	fmt.Println(mode)
	switch mode {
	case entity.CappingNone:
		return noCappingMode(ads)
	case entity.CappingReset:
		return resetCappingMode(copID, ads, ep)
	case entity.CappingStrict:
		return strictCappingMode(copID, ads, ep)
	}
	panic("invalid capping mode")
}

func strictCappingMode(copID string, ads []entity.Creative, ep string) []entity.Creative {
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

	resp := make([]entity.Creative, 0, len(ads))
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

	return []entity.Creative(resp)
}

// GetCapping try to get capping for current ad
func resetCappingMode(copID string, ads []entity.Creative, ep string) []entity.Creative {
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
		entity.Creative
	})
	resp := make([]entity.Creative, 0, len(ads))
	for i := range ads {
		size := ads[i].Size()
		if _, ok := has[size]; ok {
			has[size] = 0
		}

		key := fmt.Sprintf(
			"%s_%d",
			adKey,
			ads[i].ID(),
		)

		view := results[key]
		n := float64(view) / float64(ads[i].Campaign().Frequency())
		if n <= 1 && !selected[ads[i].ID()] {
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
			has[size]++

		} else if n > 1 {
			// capping is passed
			done[size] = append(done[size], struct {
				Key  string
				View int64
				entity.Creative
			}{
				Key:      key,
				View:     view,
				Creative: ads[i],
			})
		}
	}

	for size := range has {
		if has[size] > 0 {
			// we have one campaign with lesser capping, no need to reset, but since we are in relaxed capping mode(aka reset)
			// add over capped to pool again
			for i := range done[size] {
				capp := NewCapping(
					c,
					done[size][i].Campaign().ID(),
					done[size][i].Campaign().Frequency(),
					entity.CappingReset,
					copID,
				)
				capp.IncView(ads[i].ID(), int(done[size][i].View), false)
				done[size][i].SetCapping(capp)
				resp = append(resp, done[size][i].Creative)
			}
		} else {
			// reset this size
			var sizedCap = make([]string, len(done[size]))
			for i := range done[size] {
				capp := NewCapping(
					c,
					done[size][i].Campaign().ID(),
					done[size][i].Campaign().Frequency(),
					entity.CappingReset,
					copID,
				)
				sizedCap[i] = done[size][i].Key
				done[size][i].SetCapping(capp)
				resp = append(resp, done[size][i].Creative)
			}
			logrus.Debugf("remove key for size %d", size)
			_ = ck.Drop(sizedCap...)
		}
	}
	return []entity.Creative(resp)
}
