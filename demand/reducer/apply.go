package reducer

import (
	"context"

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

	for _, f := range ff {
		var fads = make(map[int64]entity.Creative)
		var err error
		for i := range ads {
			if ferr := f.Check(imp, ads[i]); ferr != nil {
				err = ferr
			}
		}
		if len(fads) == 0 && len(ads) != 0 {
			if imp.Publisher().Supplier().Name() != "clickyab" {
				iqs := kv.NewAEAVStore(fmt.Sprintf("DEQS_%s", time.Now().Truncate(time.Hour*24).Format("060102")), time.Hour*72)
				iqs.IncSubKey(fmt.Sprintf("%s_%s_%s", imp.Publisher().Supplier().Name(), time.Now().Truncate(time.Hour).Format("15"), err.Error()), 1)
			}
			return nil
		}
		mads = append(mads, fads)
	}

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

	return m
}
