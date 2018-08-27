package reducer

import (
	"context"

	"sync"

	"fmt"

	"time"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/kv"
)

// Filter is the interface to filter ads
type Filter interface {
	Check(entity.Context, entity.Creative) error
}

// Apply get the data and then call filter on each of them concurrently, the
// result is the accepted items
func Apply(_ context.Context, imp entity.Context, ads []entity.Creative, ff []Filter) []entity.Creative {
	var mads = make([]map[int64]entity.Creative, 0)

	var m = make([]entity.Creative, 0, len(ads))
	lock := sync.RWMutex{}
	var done = make(chan int)
	var cancel = make(chan error)

	for _, f := range ff {
		go func(f Filter) {
			var fads = make(map[int64]entity.Creative)
			var err error
			for i := range ads {
				if ferr := f.Check(imp, ads[i]); ferr != nil {
					err = ferr
					continue
				}
				fads[ads[i].ID()] = ads[i]
			}
			if len(fads) == 0 {
				cancel <- err
				close(cancel)
			}
			lock.Lock()
			mads = append(mads, fads)
			lock.Unlock()
			done <- 0
		}(f)
	}

	var c int
SELECT:
	for {
		select {
		case <-done:
			c++
			if len(ff) == c {
				ref := mads[0]
			LOOP:
				for kr, vr := range ref {
					for _, v := range mads[1:] {
						if _, ok := v[kr]; !ok {
							continue LOOP
						}
					}

					m = append(m, vr)
				}
				break SELECT
			}

		case err := <-cancel:
			iqs := kv.NewAEAVStore(fmt.Sprintf("DEQS_%s", time.Now().Truncate(time.Hour*24).Format("060102")), time.Hour*72)
			iqs.IncSubKey(fmt.Sprintf("%s_%s_%s", imp.Publisher().Supplier().Name(), time.Now().Truncate(time.Hour).Format("15"), err.Error()), 1)
			return nil
		}
	}

	return m
}
