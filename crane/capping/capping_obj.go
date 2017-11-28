package capping

import (
	"clickyab.com/crane/crane/entity"
)

// Capping is the structure for capping
type capping struct {
	view      int
	ads       map[int64]int
	sizes     map[int64]int
	frequency int
	selected  bool
}

// context is the type used to handle capping locker
type context map[int64]entity.Capping

// NewCapping create new capping
func (caps context) NewCapping(cpID int64, view, freq int) entity.Capping {
	if _, ok := caps[cpID]; !ok {
		caps[cpID] = &capping{
			view:      view,
			frequency: freq,
			ads:       make(map[int64]int),
			sizes:     make(map[int64]int),
		}
	}

	return caps[cpID]
}

func (c *capping) View() int {
	return c.view
}

func (c *capping) AdView(ad int64) int {
	return c.ads[ad]
}

func (c *capping) Frequency() int {
	return c.frequency
}

func (c *capping) Capping() int {
	return c.view / c.frequency
}

func (c *capping) AdCapping(ad int64) int {
	return c.ads[ad] / c.frequency
}

func (c *capping) IncView(ad int64, a int, sel bool) {
	c.view += a
	c.ads[ad] += a
	if sel {
		c.selected = true
	}
}

func (c *capping) Selected() bool {
	return c.selected
}
