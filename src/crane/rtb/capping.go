package rtb

import (
	"fmt"
	"services/eav"
	"sort"
	"strconv"
	"time"

	"context"
	"services/assert"

	"crane/entity"
)

const (
	userCapKey string = "CAP"
	adCapKey          = "AD"
)

// CappingContext is the type used to handle capping locker
type cappingContext map[int64]entity.Capping

// Capping is the structure for capping
type capping struct {
	adID      int64
	view      int
	frequency int
	selected  bool

	kiwi eav.Kiwi
}

// NewCapping create new capping
func (caps cappingContext) NewCapping(adID int64, freq int, kiwi eav.Kiwi) entity.Capping {
	if _, ok := caps[adID]; !ok {
		caps[adID] = &capping{
			adID:      adID,
			frequency: freq,
			kiwi:      kiwi,
		}
	}

	return caps[adID]
}

func (c *capping) View() int {
	return c.view
}

func (c *capping) Frequency() int {
	return c.frequency
}

func (c *capping) IncView(a int, sel bool) {
	c.view += a
	capKey := fmt.Sprintf("%s_%d", adCapKey, c.adID)
	c.kiwi.SetSubKey(capKey, fmt.Sprint(c.adID))
	if sel {
		c.selected = true
	}
}

func (c *capping) Selected() bool {
	return c.selected
}

func getCappingKey(copID int64) string {
	return fmt.Sprintf(
		"%s_%d_%s",
		userCapKey,
		copID,
		time.Now().Format("060102"),
	)
}

func getCapping(ctx context.Context, clientID int64, ads map[int][]entity.Advertise, slots []entity.Slot) map[int][]entity.Advertise {
	kiwi := eav.NewEavStore(getCappingKey(clientID))
	go func() {
		c := ctx.Done()
		assert.NotNil(c)
		// the channel is going to close anyway
		<-c
		// The key is hardcoded with the today (as in YYYYMMDD), so no need to config this
		kiwi.Save(24 * time.Hour)
	}()
	c := make(cappingContext)
	caps := kiwi.AllKeys()
	capsInt := make(map[string]int64)
	for s := range slots {
		size := slots[s].Size()
		// First check for ads one by one
		for a := range ads[size] {
			capKey := fmt.Sprintf("%s_%d", adCapKey, ads[size][a].ID())
			// check if the cap for this ad is full
			view, _ := strconv.ParseInt(caps[capKey], 10, 0)
			capsInt[capKey] = view
			capp := c.NewCapping(
				ads[size][a].Campaign().ID(),
				ads[size][a].Campaign().Frequency(),
				kiwi,
			)
			capp.IncView(int(view), false)
		}

		sortCap := entity.SortByCap(ads[size])
		sort.Sort(sortCap)
		ads[size] = []entity.Advertise(sortCap)
	}

	return ads
}
