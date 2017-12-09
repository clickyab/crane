package rtb

import (
	"context"

	"clickyab.com/crane/crane/builder"
	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/models"
	"clickyab.com/crane/crane/reducer"
)

// Select is the main entry point for this module
func Select(c context.Context, sel reducer.FilterFunc, opt ...builder.ShowOptionSetter) (entity.Context, error) {
	// Build context
	ctx, err := builder.NewContext(opt...)
	if err != nil {
		return nil, err
	}
	// Apply filters
	// TODO : after selector fix it
	ads := reducer.Apply(c, ctx, models.GetAds(), sel)

	// select ads
	selectAds(c, ctx, ads)

	return ctx, nil
}
