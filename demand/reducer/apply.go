package reducer

import (
	"context"

	"clickyab.com/crane/demand/entity"
)

// Filter is the interface to filter ads
type Filter interface {
	Check(entity.Context, entity.Advertise) bool
}

type mixer struct {
	f []Filter
}

func (m *mixer) Check(c entity.Context, a entity.Advertise) (b bool) {
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
func Apply(ctx context.Context, imp entity.Context, ads []entity.Advertise, ff Filter) map[int][]entity.Advertise {
	m := make(map[int][]entity.Advertise)
	for i := range ads {
		if ff.Check(imp, ads[i]) {
			n := ads[i]
			// TODO : There is a problem here. if the size is allowed in more other size then what?
			key := n.Size()
			m[key] = append(m[key], n)
		}
	}
	return m
}