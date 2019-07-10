package nrtb

import (
	"context"

	"clickyab.com/crane/models/campaign"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/reducer"
)

// Select is the main entry point for this module
func Select(c context.Context, sel []reducer.Filter, opt ...builder.ShowOptionSetter) (entity.Context, error) {
	// Build context
	ctx, err := builder.NewContext(opt...)
	if err != nil {
		return nil, err
	}
	// Apply filters
	// TODO : after selector fix it
	all, err := reducer.Apply(c, ctx, campaign.GetCampaigns(), sel)
	if err != nil {
		return nil, err
	}

	// select all
	selectAds(ctx, all)

	return ctx, nil
}
