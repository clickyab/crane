package ads

import (
	"clickyab.com/crane/models/internal/entities"
)

// AddConversion insert to conversion table
func AddConversion(wID, appID, wpID, caID, adID, copID, cpID, slotID, impID int64, acid string) error {
	return entities.AddConversion(wID, appID, wpID, caID, adID, copID, cpID, slotID, impID, acid)
}
