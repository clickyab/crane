package rtb

import (
	"context"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/reducer"
	"clickyab.com/crane/models/ads"
)

// Select is the main entry point for this module
func Select(c context.Context, sel reducer.Filter, opt ...builder.ShowOptionSetter) (entity.Context, error) {
	// Build context
	ctx, err := builder.NewContext(opt...)
	if err != nil {
		return nil, err
	}
	// Apply filters
	// TODO : after selector fix it
	all := reducer.Apply(c, ctx, ads.GetAds(), sel)
	// select all
	selectAds(c, ctx, all)

	return ctx, nil
}
