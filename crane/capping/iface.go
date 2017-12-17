package capping

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"clickyab.com/crane/crane/entity"
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

func getCappingKey(copID string) string {
	return fmt.Sprintf(
		"%s_%s_%s",
		cappingKey,
		copID,
		time.Now().Format("060102"),
	)
}

// EmptyCapping is a hack to handle no capping situation
func EmptyCapping(ads map[int][]entity.Advertise) map[int][]entity.Advertise {
	c := make(context)
	for i := range ads {
		for j := range ads[i] {
			capp := c.NewCapping(
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

// GetCapping try to get capping for current ad
func GetCapping(copID string, ads map[int][]entity.Advertise, ep string, slots ...entity.Seat) map[int][]entity.Advertise {
	var selected = make(map[int64]bool)
	if ep != "" {
		s := kv.NewDistributedSet(ep)

		for _, v := range s.Members() {
			vInt, _ := strconv.ParseInt(v, 10, 0)
			selected[vInt] = true
		}

	}
	c := make(context)
	ck := kv.NewAEAVStore(getCappingKey(copID), dailyCapExpire.Duration())
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
			n := int(view) / ads[size][ad].Campaign().Frequency()
			if n <= 1 {
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
			capp := c.NewCapping(
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
func StoreCapping(copID string, adID int64) int64 {
	return kv.NewAEAVStore(getCappingKey(copID), dailyCapExpire.Duration()).IncSubKey(fmt.Sprintf("%s_%d", adKey, adID), 1)
}
