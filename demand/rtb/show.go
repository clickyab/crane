package rtb

import (
	"context"
	"sort"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/capping"
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/reducer"
	"clickyab.com/crane/models/ads"
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
	all, err := reducer.Apply(c, ctx, ads.GetAds(), sel)
	if err != nil {
		return nil, err
	}
	// select all
	selectAds(c, ctx, all)

	return ctx, nil
}

// MinimalSelect return all ads sorted
func MinimalSelect(
	c context.Context, ctx *builder.Context, seat entity.Seat, all []entity.Creative) (
	[]entity.SelectedCreative, []entity.SelectedCreative) {
	exceed, underfloor := selector(ctx, all, seat, false, nil)

	exceed = capping.ApplyCapping(ctx.Capping(), ctx.User().ID(), exceed, ctx.EventPage())
	underfloor = capping.ApplyCapping(ctx.Capping(), ctx.User().ID(), underfloor, ctx.EventPage())

	ef := byMulti{
		Ads:   exceed,
		Video: false,
	}
	uf := byMulti{
		Ads:   underfloor,
		Video: false,
	}
	sort.Sort(ef)
	sort.Sort(uf)

	return ef.Ads, uf.Ads
}
