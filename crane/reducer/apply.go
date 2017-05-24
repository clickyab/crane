package reducer

import (
	"context"

	"clickyab.com/exchange/crane/entity"
)

var videoSize = []int{3, 4, 9, 16, 14}

// FilterFunc is the type use to filter the
type FilterFunc func(entity.Impression, entity.Advertise) bool

// Mix try to mix multiple filter to single function so there is no need to
// call Apply more than once
func Mix(f ...FilterFunc) FilterFunc {
	return func(e entity.Impression, a entity.Advertise) bool {
		for i := range f {
			if !f[i](e, a) {
				return false
			}
		}
		return true
	}
}

func checkCtx(d <-chan struct{}) bool {
	if d == nil {
		return true
	}
	select {
	case <-d:
		return false
	default:
		return true
	}
}

// Apply get the data and then call filter on each of them concurrently, the
// result is the accepted items
func Apply(ctx context.Context, imp entity.Impression, ads []entity.Advertise, ff FilterFunc) map[int][]entity.Advertise {
	d := ctx.Done()
	m := make(map[int][]entity.Advertise)
	for i := range ads {
		if !checkCtx(d) {
			break
		}
		if ff(imp, ads[i]) {
			n := ads[i].Copy()
			n.SetWinnerBID(0)
			if n.Type() == entity.AdTypeVideo {
				for _, j := range videoSize {
					m[j] = append(m[j], n)
				}
			} else {
				m[ads[i].Size()] = append(m[ads[i].Size()], n)
			}
		}
	}
	return m
}
