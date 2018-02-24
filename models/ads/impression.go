package ads

import (
	"time"

	"clickyab.com/crane/models/internal/entities"
)

// FindImpressionByID return impression by impression id
func FindImpressionByID(impID int64, t time.Time) (*entities.Impression, error) {
	return entities.FindImpressionByID(impID, t)
}

// FindImpressionByRH return impression by reserved hash
func FindImpressionByRH(rh string, t time.Time) (*entities.Impression, error) {
	return entities.FindImpressionByRH(rh, t)
}

// FindImpFromClickByImpID return impression by impression id
func FindImpFromClickByImpID(impID int64) (*entities.Impression, error) {
	return entities.FindImpFromClickByImpID(impID)
}

// FindImpFromClickByRH return impression by reserved hash
func FindImpFromClickByRH(rh string) (*entities.Impression, error) {
	return entities.FindImpFromClickByRH(rh)
}
