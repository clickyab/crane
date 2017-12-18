package rtb

import (
	"context"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/models"
	"clickyab.com/crane/demand/reducer"
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
	ads := reducer.Apply(c, ctx, models.GetAds(), sel)
	// select ads
	selectAds(c, ctx, ads)

	return ctx, nil
}
