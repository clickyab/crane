package capping

import (
	"sync"

	"fmt"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/kv"
)

// Capping is the structure for capping
type capping struct {
	view      int
	lock      sync.RWMutex
	ads       map[int64]int
	frequency int
	selected  bool
	mode      entity.CappingMode
	copID     string
}

// Context is for capping context
type Context struct {
	m map[int64]entity.Capping
	l *sync.Mutex
}

func newContext() *Context {
	return &Context{
		m: make(map[int64]entity.Capping),
		l: &sync.Mutex{},
	}
}

// NewCapping create new capping
func NewCapping(ctx *Context, cpID int64, freq int, mode entity.CappingMode, uid string) entity.Capping {
	ctx.l.Lock()
	defer ctx.l.Unlock()

	if _, ok := ctx.m[cpID]; !ok {
		ctx.m[cpID] = &capping{
			frequency: freq,
			ads:       make(map[int64]int),
			mode:      mode,
			copID:     uid,
		}
	}

	return ctx.m[cpID]
}

func (c *capping) View() int {
	return c.view
}

func (c *capping) AdView(ad int64) int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	v, ok := c.ads[ad]
	if ok {
		return v
	}
	return 0
}

func (c *capping) Frequency() int {
	return c.frequency
}

func (c *capping) Capping() int {
	return c.view / c.frequency
}

func (c *capping) Store(ad int64) {
	if c.mode == entity.CappingNone {
		return
	}
	kv.NewAEAVStore(getCappingKey(c.mode, c.copID), dailyCapExpire.Duration()).IncSubKey(fmt.Sprintf("%s_%d", adKey, ad), 1)
}

func (c *capping) AdCapping(ad int64) int {
	return c.AdView(ad) / c.frequency
}

func (c *capping) IncView(ad int64, a int, sel bool) {
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
