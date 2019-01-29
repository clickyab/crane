package campaign

import (
	"fmt"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/slack"
)

// Strategy checker
type Strategy struct {
}

// Check check if creative can be use for this impression
func (*Strategy) Check(impression entity.Context, ad entity.Campaign) error {
	// TODO :// just for debugging
	if ad == nil {
		go func() {
			slack.AddCustomSlack(fmt.Errorf("[WTF]campaign is null for following ad id %d", ad.ID()))
		}()
	}
	if ad.Strategy().IsSubsetOf(impression.Strategy()) {
		return nil
	}
	return fmt.Errorf("supplier strategy is %d but campaign want %d ",
		impression.Strategy(), ad.Strategy())
}
