package filter

import (
	"clickyab.com/gad/models"
	"clickyab.com/crane/crane/builder"
)

// RemoveSlotPins remove fix slot from ad pool
func RemoveSlotPins(c *builder.context, in models.AdData) bool {
	// TODO : revert this, after fixing the slot pin
	//for i := range c.SlotPins {
	//	if c.SlotPins[i].AdID == in.AdID {
	//		return false
	//	}
	//}
	return true
}
