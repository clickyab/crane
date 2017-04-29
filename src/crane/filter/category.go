package filter

import (
	"crane/entity"
)

// Category check for category match
func Category(impression entity.Impression, advertise entity.Advertise) bool {
	whitelist := advertise.Category()
	elems := impression.Category()

	return hasCategory(elems, whitelist)
}
