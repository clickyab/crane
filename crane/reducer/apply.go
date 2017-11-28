package reducer

import (
	"context"
	"fmt"

	"clickyab.com/crane/crane/entity"
)

// FilterFunc is the type use to filter the
type FilterFunc func(entity.Context, entity.Advertise) bool

// Mix try to mix multiple filter to single function so there is no need to
// call Apply more than once
func Mix(f ...FilterFunc) FilterFunc {
	return func(e entity.Context, a entity.Advertise) bool {
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
func Apply(ctx context.Context, imp entity.Context, ads []entity.Advertise, ff FilterFunc) map[string][]entity.Advertise {
	d := ctx.Done()
	m := make(map[string][]entity.Advertise)
	for i := range ads {
		if !checkCtx(d) {
			break
		}
		if ff(imp, ads[i]) {
			n := ads[i].Duplicate()
			n.SetWinnerBID(0)
			// TODO : There is a problem here. if the size is allowed in more other size then what?
			key := fmt.Sprintf("%dx%d", ads[i].Width(), ads[i].Height())
			m[key] = append(m[key], n)
		}
	}
	return m
}
