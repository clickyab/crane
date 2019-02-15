package capping

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/kv"
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
func ApplyCapping(mode entity.CappingMode, copID string, ads []entity.SelectedCreative) []entity.SelectedCreative {
	switch mode {
	case entity.CappingNone:
		return noCappingMode(ads)
	case entity.CappingReset:
		return resetCappingMode(copID, ads)
	case entity.CappingStrict:
		return strictCappingMode(copID, ads)
	}
	panic("invalid capping mode")
}

func strictCappingMode(copID string, ads []entity.SelectedCreative) []entity.SelectedCreative {
	var selected = make(map[int32]bool)

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

		view := int32(results[key])
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
			capp.IncView(ads[i].ID(), view, false)
			passed.SetCapping(capp)
			resp = append(resp, passed)
		}
	}

	return resp
}

// GetCapping try to get capping for current ad
func resetCappingMode(copID string, ads []entity.SelectedCreative) []entity.SelectedCreative {
	var selected = make(map[int32]bool)
	c := newContext()
	ck := kv.NewAEAVStore(getCappingKey(entity.CappingReset, copID), dailyCapExpire.Duration())
	results := ck.AllKeys()
	has := make(map[int32]int32)
	done := make(map[int32][]struct {
		Key  string
		View int32
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

		view := int32(results[key])
		n := float64(view) / float64(ads[i].Campaign().Frequency())
		if n < 1 && !selected[ads[i].ID()] {
			capp := NewCapping(
				c,
				ads[i].Campaign().ID(),
				ads[i].Campaign().Frequency(),
				entity.CappingReset,
				copID,
			)
			capp.IncView(ads[i].ID(), view, false)
			ads[i].SetCapping(capp)
			resp = append(resp, ads[i])
			has[size] = has[size] + 1
		} else if !selected[ads[i].ID()] {
			// capping is passed
			done[size] = append(done[size], struct {
				Key  string
				View int32
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
			capp.IncView(done[size][i].ID(), done[size][i].View, false)
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
