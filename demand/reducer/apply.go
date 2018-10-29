package reducer

import (
	"context"
	"errors"
	"time"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/metrics"
	"github.com/clickyab/services/xlog"
	"github.com/prometheus/client_golang/prometheus"
)

// Filter is the interface to filter ads
type Filter interface {
	Check(entity.Context, entity.Creative) error
}

type filtered struct {
	ad      entity.Creative
	confirm int
}

// ErrorTimeOut for timeout filter
var ErrorTimeOut = errors.New("timeout")

// Apply get the data and then call filter on each of them concurrently, the
// result is the accepted items
func Apply(c context.Context, imp entity.Context, ads []entity.Creative, ff []Filter) ([]entity.Creative, error) {
	var mads = make(map[int32]*filtered)
	var res = make([]entity.Creative, 0)
	fads := make(chan entity.Creative)
	fcl := make(chan error)
	done := make(chan int)
	next := make(chan bool)
	dl := time.After(time.Millisecond * 500)

	for _, f := range ff {
		go func(f Filter) {
			c := 0
			var err error
			for i := range ads {
				fe := f.Check(imp, ads[i])
				if fe == nil {
					c++
					fads <- ads[i]
					if _, ok := <-next; !ok {
						return
					}
				}
				err = fe
			}
			if c == 0 {
				fcl <- err
			} else {
				done <- 0
			}
		}(f)
	}

	var counter = 0

LOOP:
	for {
		select {
		case <-c.Done():
			err := errors.New("canceled")
			go metrics.Filter.With(
				prometheus.Labels{
					"supplier": imp.Publisher().Supplier().Name(),
					"reason":   err.Error(),
				},
			).Inc()
			close(next)
			return nil, err
		case res := <-fcl:
			xlog.Get(c).Debugf("Filter doesn't match: %s", res)
			go metrics.Filter.With(
				prometheus.Labels{
					"supplier": imp.Publisher().Supplier().Name(),
					"reason":   res.Error(),
				},
			).Inc()
			close(next)
			return nil, res
		case <-dl:
			xlog.Get(c).Debugf("Filter timeout")
			go metrics.Filter.With(
				prometheus.Labels{
					"supplier": imp.Publisher().Supplier().Name(),
					"reason":   ErrorTimeOut.Error(),
				},
			).Inc()
			close(next)
			return nil, ErrorTimeOut
		case <-done:
			counter++
			if len(ff) == counter {
				break LOOP
			}
		case ad := <-fads:
			next <- true
			if v, ok := mads[ad.ID()]; ok {
				v.confirm++
				continue
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
	return res, nil
}
