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
func noCappingMode(ads []entity.Advertise) []entity.Advertise {
	c := newContext()
	for i := range ads {
		capp := NewCapping(
			c,
			ads[i].Campaign().ID(),
			0,
			ads[i].Campaign().Frequency(),
		)
		ads[i].SetCapping(capp)

	}

	return ads
}

// ApplyCapping is an entry for set capping in ads
func ApplyCapping(mode entity.CappingMode, copID string, ads []entity.Advertise, ep string, slots ...entity.Seat) []entity.Advertise {
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

func strictCappingMode(copID string, ads []entity.Advertise, ep string, slots ...entity.Seat) []entity.Advertise {
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

	resp := make([]entity.Advertise, 0, len(ads))
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
				0,
				ads[i].Campaign().Frequency(),
			)
			capp.IncView(passed.ID(), int(view), false)
			passed.SetCapping(capp)
			resp = append(resp, passed)
		}
	}

	return resp
}

// GetCapping try to get capping for current ad
func resetCappingMode(copID string, ads []entity.Advertise, ep string, slots ...entity.Seat) []entity.Advertise {
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
		entity.Advertise
	})
	resp := make([]entity.Advertise, 0, len(ads))
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
				0,
				ads[i].Campaign().Frequency(),
			)
			capp.IncView(ads[i].ID(), int(view), false)
			ads[i].SetCapping(capp)
			resp = append(resp, ads[i])
			has[size] += 1
		} else if n > 1 {
			// capping is passed
			done[size] = append(done[size], struct {
				Key  string
				View int64
				entity.Advertise
			}{
				Key:       key,
				View:      view,
				Advertise: ads[i],
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
					0,
					done[size][i].Campaign().Frequency(),
				)
				capp.IncView(done[size][i].ID(), int(done[size][i].View), false)
				done[size][i].SetCapping(capp)
				resp = append(resp, done[size][i].Advertise)
			}
		} else {
			// reset this size
			var sizedCap = make([]string, len(done[size]))
			for i := range done[size] {
				capp := NewCapping(
					c,
					done[size][i].Campaign().ID(),
					0,
					done[size][i].Campaign().Frequency(),
				)
				sizedCap[i] = done[size][i].Key
				done[size][i].SetCapping(capp)
				resp = append(resp, done[size][i].Advertise)
			}
			logrus.Debugf("remove key for size %d", size)
			ck.Drop(sizedCap...)
		}
	}
	return ads
}

// StoreCapping try to store a capping object
func StoreCapping(mode entity.CappingMode, copID string, adID int64) int64 {
	return kv.NewAEAVStore(getCappingKey(mode, copID), dailyCapExpire.Duration()).IncSubKey(fmt.Sprintf("%s_%d", adKey, adID), 1)
}
