package rtb

import (
	"context"

	"clickyab.com/crane/crane/builder"
	"clickyab.com/crane/crane/entity"
	"clickyab.com/gad/models/selector"
)

func Select(c context.Context, sel selector.FilterFunc, opt ...builder.ShowOptionSetter) (map[string]string, map[string]entity.Advertise, error) {
	// Build context
	ctx, err := builder.NewContext(opt...)
	if err != nil {
		return nil, nil, err
	}
	// Apply filters
	// TODO : after selector fix it
	ads := map[int][]entity.Advertise{} // selector.Apply(ctx, selector.GetAdData(), sel)

	// select ads
	links, selected := selectAds(c, ctx, ads)
	return links, selected, nil
}
