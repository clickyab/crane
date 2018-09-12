package capping

import (
	"sync"

	"fmt"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/kv"
)

// Capping is the structure for capping
type capping struct {
	view      int32
	lock      sync.RWMutex
	ads       map[int32]int32
	frequency int32
	selected  bool
	mode      entity.CappingMode
	copID     string
}

// Context is for capping context
type Context struct {
	m map[int32]entity.Capping
	l *sync.Mutex
}

func newContext() *Context {
	return &Context{
		m: make(map[int32]entity.Capping),
		l: &sync.Mutex{},
	}
}

// NewCapping create new capping
func NewCapping(ctx *Context, cpID int32, freq int32, mode entity.CappingMode, uid string) entity.Capping {
	ctx.l.Lock()
	defer ctx.l.Unlock()

	if _, ok := ctx.m[cpID]; !ok {
		ctx.m[cpID] = &capping{
			frequency: freq,
			ads:       make(map[int32]int32),
			mode:      mode,
			copID:     uid,
		}
	}

	return ctx.m[cpID]
}

func (c *capping) View() int32 {
	return c.view
}

func (c *capping) AdView(ad int32) int32 {
	c.lock.RLock()
	defer c.lock.RUnlock()

	v, ok := c.ads[ad]
	if ok {
		return v
	}
	return 0
}

func (c *capping) Frequency() int32 {
	return c.frequency
}

func (c *capping) Capping() int32 {
	return c.view / c.frequency
}

func (c *capping) AdCapping(ad int32) int32 {
	return c.AdView(ad) / c.frequency
}

func (c *capping) IncView(ad int32, a int32, sel bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.view += a
	c.ads[ad] += a
	if sel {
		c.selected = true
	}
}

func (c *capping) Selected() bool {
	return c.selected
}

// StoreCapping try to store a capping object
func StoreCapping(mode entity.CappingMode, copID string, adID int64) int64 {
	return kv.NewAEAVStore(getCappingKey(mode, copID), dailyCapExpire.Duration()).IncSubKey(fmt.Sprintf("%s_%d", adKey, adID), 1)
}
