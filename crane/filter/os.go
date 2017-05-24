package filter

import (
	"clickyab.com/exchange/crane/entity"
)

// OS check for os matched in filters
func OS(impression entity.Impression, advertise entity.Advertise) bool {
	blacklist := advertise.AllowedOS()
	if len(blacklist) == 0 {
		// No os is blacklisted, so pass it
		return true
	}

	elem := impression.OS().ID

	return hasInt64(blacklist, elem)
}
