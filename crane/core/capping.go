package core

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/kv"
)

const (
	prefix = `cap_`
)

type cap struct {
	ads []entity.Advertise
	kh  kv.AKiwi
	kt  kv.AKiwi
}

// Capping interface capping object for all the campaign ads
type Capping interface {
	Sort() []entity.Advertise
	IncView(advertise entity.Advertise)
}

// GetCapping return new cap
func GetCapping(clientID string, ads []entity.Advertise) Capping {
	kh := kv.NewAEAVStore(fmt.Sprintf("%s_%s_%s", prefix, clientID, time.Now().Format("15")))
	kt := kv.NewAEAVStore(fmt.Sprintf("%s_%s_%s", prefix, clientID, time.Now().Format("0102")))
	// to cache all sub keys

	return &cap{
		kh:  kh,
		kt:  kt,
		ads: ads,
	}
}

// Get sorted ads
func (d *cap) Sort() []entity.Advertise {
	sort.Sort(d)
	return d.ads
}

// IncView add a view to counter
func (d *cap) IncView(a entity.Advertise) {
	n := time.Now()
	k := strconv.FormatInt(a.Campaign().ID(), 64)

	d.kh.IncSubKey(k, 1)
	hx, err := time.Parse("06010215", n.Format("06010215"))
	assert.Nil(err)
	err = d.kh.Save(hx.Add(time.Hour).Sub(n))
	assert.Nil(err)
	tx, err := time.Parse("060102", n.Format("060102"))
	assert.Nil(err)
	d.kt.IncSubKey(k, 1)
	err = d.kt.Save(tx.AddDate(0, 0, 1).Sub(n))
	assert.Nil(err)

}

// Len for sorting
func (d *cap) Len() int {
	return len(d.ads)
}

// Swap for sorting
func (d *cap) Swap(i, j int) {
	d.ads[i], d.ads[j] = d.ads[j], d.ads[i]
}

// Less for sorting
func (d *cap) Less(i, j int) bool {

	ka := strconv.FormatInt(d.ads[i].Campaign().ID(), 64)
	kb := strconv.FormatInt(d.ads[j].Campaign().ID(), 64)

	ha := d.kh.AllKeys()[ka]
	hb := d.kh.AllKeys()[kb]
	ta := d.kt.AllKeys()[ka]
	tb := d.kt.AllKeys()[kb]

	fa := isFull(d.ads[i].Campaign().Frequency() <= int(ta))
	fb := isFull(d.ads[j].Campaign().Frequency() <= int(tb))
	ca, cb := checkCpm(d.ads[i].CPM(), d.ads[i].CPM())

	return fmt.Sprintf("%s%03x%03x%d", fa, ha, ta, ca) < fmt.Sprintf("%s%03x%03x%d", fb, hb, tb, cb)

}
func checkCpm(a, b int64) (ra, rb int) {
	if a == b {
		return 0, 0
	}
	if a < b {
		return 1, 0
	}
	return 0, 1
}

type isFull bool

// String for stringer
func (f isFull) String() string {
	if f {
		return "1"
	}
	return "0"
}
