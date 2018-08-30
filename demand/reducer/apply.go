package reducer

import (
	"context"

	"time"

	"clickyab.com/crane/demand/entity"
)

// Filter is the interface to filter ads
type Filter interface {
	Check(entity.Context, entity.Creative) error
}

type filtered struct {
	ad      entity.Creative
	confirm int
}

// Apply get the data and then call filter on each of them concurrently, the
// result is the accepted items
func Apply(c context.Context, imp entity.Context, ads []entity.Creative, ff []Filter) []entity.Creative {
	var mads = make(map[int64]*filtered)
	var res = make([]entity.Creative, 0)
	fads := make(chan entity.Creative)
	ctx, cl := context.WithCancel(c)
	dl := time.After(time.Millisecond * 20)
	for _, f := range ff {
		go func(f Filter) {
			c := 0
			for i := range ads {
				if ferr := f.Check(imp, ads[i]); ferr == nil {
					c++
					fads <- ads[i]
				}
			}
			if c == 0 {
				cl()
			}
		}(f)
	}

LOOP:
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-dl:
			return nil
		case ad := <-fads:
			if v, ok := mads[ad.ID()]; ok {
				v.confirm++
				if v.confirm == len(ff) {
					break LOOP
				}
			}
			mads[ad.ID()] = &filtered{
				ad:      ad,
				confirm: 1,
			}
		}
	}

	for _, v := range mads {
		if v.confirm == len(ff) {
			res = append(res, v.ad)
		}
	}
	return res
}
