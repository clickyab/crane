package filter

import (
	"errors"

	"clickyab.com/crane/demand/entity"
)

// Strategy checker
type Strategy struct {
}

// Check check if creative can be use for this impression
func (*Strategy) Check(impression entity.Context, ad entity.Creative) error {
	if ad.Campaign().Strategy().IsSubsetOf(impression.Strategy()) {
		return nil
	}
	return errors.New("strategy filter not met")

}
