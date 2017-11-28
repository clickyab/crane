package filter

import "clickyab.com/crane/crane/entity"

// OS check for os matched in filters
func OS(impression entity.Context, advertise entity.Advertise) bool {
	blacklist := advertise.Campaign().AllowedOS()
	if len(blacklist) == 0 {
		// No os is blacklisted, so pass it
		return true
	}

	elem := impression.OS().Name

	return hasString(blacklist, elem)
}
