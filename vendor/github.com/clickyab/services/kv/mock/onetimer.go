package mock

import (
	"sync"
	"time"
)

var (
	onetimes = make(map[string]string)
	otLock   sync.Mutex

	_ oneTimer
)

type oneTimer struct {
	d   time.Duration
	key string
}

func (ot *oneTimer) Key() string {
	return ot.key
}

func (ot *oneTimer) Set(s string) string {
	otLock.Lock()
	defer otLock.Unlock()

	if v, ok := onetimes[ot.key]; ok {
		return v
	}

	onetimes[ot.key] = s
	return s
}
