package mock

import (
	"sync"
	"time"

	"github.com/clickyab/services/dset"
	"github.com/clickyab/services/safe"
)

var (
	distributedSets = make(map[string]dset.DistributedSet)
	lock            = sync.Mutex{}
)

// NewMockDsetStore retrieve DistributedSet from store or make new one if dose not exist
func NewMockDsetStore(key string) dset.DistributedSet {
	lock.Lock()
	defer lock.Unlock()
	if distributedSets[key] != nil {
		return distributedSets[key]
	}
	d := &distributedSet{
		key: key,
	}
	return d
}

type distributedSet struct {
	once    sync.Once
	members []string
	adds    []string
	key     string
	exp     time.Time
}

// Members return ads ID
func (d *distributedSet) Members() []string {
	return append(d.members, d.adds...)
}

// Add new ad ID to memebers (after invoking save)
func (d *distributedSet) Add(newMembers ...string) {
	d.adds = append(d.adds, newMembers...)
}

// Key of DistributedSet
func (d *distributedSet) Key() string {
	return d.key
}

// Save added IDs and extend TTL of distributedSet
func (d *distributedSet) Save(t time.Duration) error {
	lock.Lock()
	defer lock.Unlock()
	d.exp = time.Now().Add(t)
	d.members = append(d.members, d.adds...)
	d.adds = nil
	distributedSets[d.key] = d
	d.once.Do(func() {
		safe.GoRoutine(d.ttl)
	})
	return nil
}

func (d *distributedSet) ttl() {
	for {
		<-time.After(time.Until(d.exp))
		if time.Now().Unix() > d.exp.Unix() {
			lock.Lock()
			defer lock.Unlock()
			delete(distributedSets, d.key)
			return
		}
	}
}
