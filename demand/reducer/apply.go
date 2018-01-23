package reducer

import (
	"context"

	"clickyab.com/crane/demand/entity"
)

// Filter is the interface to filter ads
type Filter interface {
	Check(entity.Context, entity.Creative) bool
}

type mixer struct {
	f []Filter
}

func (m *mixer) Check(c entity.Context, a entity.Creative) (b bool) {
	for i := range m.f {
		if !m.f[i].Check(c, a) {
			// Un-coment this for debug
			//fmt.Printf("\nfalse on %T", m.f[i])
			return false
		}
	}
	//fmt.Printf("\true on %d", a.ID())
	return true
}

// Mix try to mix multiple filter to single function so there is no need to
// call Apply more than once
func Mix(f ...Filter) Filter {
	return &mixer{f: f}
}

// Apply get the data and then call filter on each of them concurrently, the
// result is the accepted items
func Apply(_ context.Context, imp entity.Context, ads []entity.Creative, ff Filter) []entity.Creative {
	var m = make([]entity.Creative, 0, len(ads))
	for i := range ads {
		if ff.Check(imp, ads[i]) {
			m = append(m)
		}
	}
	return m
}
